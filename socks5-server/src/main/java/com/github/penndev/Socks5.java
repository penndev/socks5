package com.github.penndev;

import java.net.ServerSocket;
import java.net.Socket;

/**
 * This is SOCKS Protocol Version 5 Java implement
 * <a href="https://datatracker.ietf.org/doc/html/rfc1928">socks5 rfc</a>
 * <a href="https://datatracker.ietf.org/doc/html/rfc1929">socks5 rfc auth</a>
 */
public class Socks5 {
    public static int Port = 1080;

    public static String username = null;

    public static String password = null;


    /**
     * Listen port and run socks5 tcp service
     *
     * @param args [port] [username] [password]
     */
    public static void main(String[] args) {
        // 获取传递的参数
        if (args.length == 1 || args.length == 3) {
            Socks5.Port = Integer.parseInt(args[0]);
            if (args.length == 3) {
                Socks5.username = args[1];
                Socks5.password = args[2];
            }
        } else if (args.length > 1) {
            System.out.println("Socks5 args:[port] [?username] [?password]");
            return;
        }
        // 打印启动参数
        System.out.printf("Socks5 listening args:[port=%d] [username=%s] [password=%s]\n",
                Socks5.Port, Socks5.username, Socks5.password);
        try (ServerSocket srvSocket = new ServerSocket(Socks5.Port)) {
            //noinspection InfiniteLoopStatement
            while (true) {
                Socket sock = srvSocket.accept();
                //System.out.printf("Conn from: [%s] \n", sock.getRemoteSocketAddress());
                Thread t = new Thread(new Service(sock));
                t.start();
            }
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }

    }
}