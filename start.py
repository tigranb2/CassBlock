from sys import modules, argv
from os import system
from functools import partial
from time import time, sleep

from mininet.net import Mininet
from mininet.node import CPULimitedHost
from mininet.link import TCLink
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

from topos import *
from cass import *
from config import conf
from geth import *

node_count = int(argv[1])
row_count = int(argv[2])
rate_parameter = int(argv[3])


def get_topology():
    # privateDirs = [('~/.ethereum', '~/%(name)s/.ethereum')]
    privateDirs = []
    host = partial(CPULimitedHost, privateDirs=privateDirs)
    try:
        topo_cls = getattr(modules[__name__], conf["topo"]["class"])
        topo_obj = topo_cls(*conf['topo']["args"], **conf['topo']["kwargs"])
        net = Mininet(topo=topo_obj, host=host, link=TCLink)
        return topo_obj, net
    except Exception as e:
        print("Specified topology not found: ", e)
        exit(0)


def test_topology(topo: Topo, net: Mininet):
    print("Dumping host connections")
    dumpNodeConnections(net.hosts)
    print("Waiting switch connections")
    net.waitConnected()

    print("Testing network connectivity - (i: switches are learning)")
    net.pingAll()
    print("Testing network connectivity - (ii: after learning)")
    net.pingAll()

    print("Get all hosts")
    print(topo.hosts(sort=True))

    # print("Get all links")
    # for link in topo.links(sort=True, withKeys=True, withInfo=True):
    #     pprint(link)
    # print()

    if conf['test']['iperf'] == -1:
        return
    else:
        hosts = [net.get(i) for i in topo.hosts(sort=True)]
        if conf['test']['iperf'] == 0:
            net.iperf((hosts[0], hosts[-1]))
        else:
            [net.iperf((i, j)) for i in hosts for j in hosts if i != j]


def main():
    def delay_command(host, cmd, print=True):
        sleep(0.5)
        if print:
            hs[host - 1].cmdPrint(cmd)
        else:
            hs[host - 1].cmd(cmd)
        sleep(0.5)

    system('sudo mn --clean')
    setLogLevel('info')

    # reads YAML configs and creates the network
    topo, net = get_topology()
    net.start()

    hs = topo.hosts(sort=True)
    hs = [net.getNodeByName(h) for h in hs]

    # adds nodes to cluster
    system("ccm add node1 -i 10.0.0.1 -j 7100 -s")
    for i in range(2, node_count + 1):
        cmd = f"ccm add node{i} -i 10.0.0.{i} -j 7{i}00 -s"
        system(cmd)

    # starts nodes
    for i in range(1, node_count + 1):
        cmd = f"ccm node{i} start"
        delay_command(i, cmd)
        sleep(3)

    delay_command(1, "ccm node1 ring")

    delay_command(1, init_keyspace)
    delay_command(1, init_table)

    # starts geth
    delay_command(1, miner_start)
    for i in range(2, node_count + 1):
        delay_command(i, gen_enode)
        delay_command(i, node_start.format(n=i))

    sleep(60)

    # starts writes
    for i in range(1, node_count):
        cmd = f"./simulateWrites {i} {row_count} {rate_parameter} &"
        delay_command(i, cmd)

    cmd = f"./simulateWrites {node_count} {row_count} {rate_parameter}"
    delay_command(node_count, cmd)

    # stop the network
    net.stop()


if __name__ == '__main__':
    main()
