from threading import Thread
import context
import client
import simulator
# import animate


def simulate(sim: simulator.Simulator):
    sim.run()


if __name__ == '__main__':
    context.read_buckets()
    context.read_contexts()

    kontext = context.banner_contexts.contexts[0]
    client = client.Client(context=kontext, host="127.0.0.1", port=8000)

    simulator = simulator.Simulator(cln=client, ctx=kontext)

    # thread = Thread(target=simulate, args=(simulator,))
    # thread.start()
    # animate.run_animate(host="127.0.0.1", port=8000, context=kontext)
    # simulator.stop()
    # thread.join()
