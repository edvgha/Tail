package space

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func (s *Space) BackgroundTask() {
	go func() {
		for {
			if s.ExplorationQty.Load()-s.LastUpdateQty.Load() >= 100 {
				s.Learn()
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

type Estimation struct {
	price float64
	pr    float64
}

type Estimations []Estimation

func (s *Space) Learn() {
	estimations := make([]Estimations, len(s.Levels))
	s.mutex.Lock()
	for i := 0; i < len(s.Levels); i++ {
		estimations[i] = make([]Estimation, len(s.Levels[i].Buckets))
		for j := 0; j < len(s.Levels[i].Buckets); j++ {
			estimations[i][j] = Estimation{
				price: s.Levels[i].Buckets[j].Lhs + (s.Levels[i].Buckets[j].Rhs-s.Levels[i].Buckets[j].Lhs)/2.0,
				pr:    s.Levels[i].Buckets[j].Pr,
			}
		}
	}
	s.mutex.Unlock()
	for i, estimation := range estimations {
		// TODO run in parallel
		// TODO or in once for all
		prs, err := learnNonDecreasing(estimation)
		if err != nil {
			log.Error().Msgf("Faild to learn for level %d", i)
		}
		for i := 0; i < len(estimation); i++ {
			estimation[i].pr = prs[i]
		}
	}

	s.wcMutex.Lock()
	for i := 0; i < len(s.Levels); i++ {
		for j := 0; j < len(s.Levels[i].Buckets); j++ {
			s.Levels[i].WinningCurve[j] = estimations[i][j].pr
			if estimations[i][j].price < s.Levels[i].Buckets[j].Lhs || estimations[i][j].price > s.Levels[i].Buckets[j].Rhs {
				log.Fatal().Msgf("inconsistency price %f in range [%f, %f]",
					estimations[i][j].price, s.Levels[i].Buckets[j].Lhs, s.Levels[i].Buckets[j].Rhs)
			}
		}
	}
	s.wcMutex.Unlock()
}

func learnNonDecreasing(estimations Estimations) ([]float64, error) {
	in, err := writeToCSV(estimations)
	defer os.Remove(in.Name())
	if err != nil {
		return nil, err
	}

	return runLightGBM(in)
}

func writeToCSV(estimations Estimations) (*os.File, error) {
	f, err := os.CreateTemp("", "regression.train")
	if err != nil {
		log.Error().Msg("failed to create tmp input data csv file")
		return nil, err
	}

	w := csv.NewWriter(f)
	w.Comma = '\t'
	estimationsStr := make([][]string, len(estimations))
	for i := 0; i < len(estimations); i++ {
		estimationsStr[i] = make([]string, 2)
		estimationsStr[i][0] = fmt.Sprintf("%f", estimations[i].pr)
		estimationsStr[i][1] = fmt.Sprintf("%f", estimations[i].price)
	}
	w.WriteAll(estimationsStr)
	if err := w.Error(); err != nil {
		log.Error().Msg("failed to write input csv data")
		return nil, err
	}
	return f, nil
}

func runLightGBM(inputData *os.File) ([]float64, error) {
	modelFile := "/tmp/LightGBM_model.txt"
	outputData := "/tmp/output.txt"
	trainConf, _, err := generateTrainConfig(inputData.Name(), modelFile)
	defer os.Remove(trainConf.Name())
	if err != nil {
		return nil, err
	}
	predictConf, _, err := generatePredictConfig(inputData.Name(), outputData, modelFile)
	defer os.Remove(predictConf.Name())
	if err != nil {
		return nil, err
	}

	prs, err := execute(trainConf.Name(), predictConf.Name(), outputData)
	os.Remove(modelFile)
	os.Remove(outputData)
	if err != nil {
		return nil, err
	}
	return prs, err
}

func generateTrainConfig(inData, modelFile string) (*os.File, int, error) {
	f, err := os.CreateTemp("", "train.conf")
	if err != nil {
		log.Error().Msg("failed to create tmp train config file")
		return nil, 0, err
	}
	defer f.Close()

	trainConf := "task = train\n" +
		"boosting_type = gbdt\n" +
		"objective = regression\n" +
		"metric = l2\n" +
		"metric_freq = 1\n" +
		"is_training_metric = false\n" +
		"label_column = 0\n" +
		"data = " + inData + "\n" +
		"num_trees = 100\n" +
		"learning_rate = 0.1\n" +
		"num_leaves = 31\n" +
		"min_child_samples = 2\n" +
		"mc = 1\n" +
		"tree_learner = serial\n" +
		"is_enable_sparse = true\n" +
		"use_two_round_loading = false\n" +
		"is_save_binary_file = false\n" +
		"output_model = " + modelFile

	n, err := f.WriteString(trainConf)
	if err != nil {
		log.Error().Msg("failed to write into train config file")
		return nil, 0, err
	}
	return f, n, nil
}

func generatePredictConfig(inData, outData, modelFile string) (*os.File, int, error) {
	f, err := os.CreateTemp("", "predict.conf")
	if err != nil {
		log.Error().Msg("failed to create tmp prediction config file")
		return nil, 0, err
	}
	defer f.Close()

	predictConf := "task = predict\n" +
		"data = " + inData + "\n" +
		"output_result = " + outData + "\n" +
		"input_model = " + modelFile
	n, err := f.WriteString(predictConf)
	if err != nil {
		log.Error().Msg("failed to write into prediction config file")
		return nil, 0, err
	}
	return f, n, nil
}

func execute(trainConf, predictConf, outData string) ([]float64, error) {
	executable, err := exec.LookPath("lightgbm")
	if err != nil {
		log.Error().Msg("failed to find 'lightgbm' executable")
		return nil, err
	}
	cmd := exec.Command(executable, "config="+trainConf)
	err = cmd.Run()
	if err != nil {
		log.Error().Msg("failed on training")
		return nil, err
	}

	cmd = exec.Command(executable, "config="+predictConf)
	err = cmd.Run()
	if err != nil {
		log.Error().Msg("failed on prediction")
		return nil, err
	}

	return readPredictions(outData)
}

func readPredictions(outData string) ([]float64, error) {
	r, err := os.Open(outData)
	if err != nil {
		log.Error().Msgf("failed to open prediction file: %s", outData)
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var prs []float64
	for scanner.Scan() {
		f64Str := scanner.Text()
		f64, err := strconv.ParseFloat(f64Str, 64)
		if err != nil {
			log.Error().Msgf("failed to parse float64 from string %s", f64Str)
			prs = append(prs, 0.0)
		} else {
			prs = append(prs, f64)
		}
	}
	return prs, nil
}
