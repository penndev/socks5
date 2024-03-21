package com.penndev.socks5.service

import android.app.PendingIntent
import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor
import android.util.Log
import android.widget.Toast
import androidx.core.app.NotificationCompat
import com.penndev.socks5.MainActivity
import com.penndev.socks5.R
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.Job
import kotlinx.coroutines.launch
import java.net.InetAddress

class Socks5ServiceCloseException(message: String) : Exception(message)

class Socks5Service : VpnService() {

    // 回调UI状态
    interface OnStatus {
        fun start() {
            status = true
        }

        fun close() {
            status = false
        }
    }

    companion object {
        var status: Boolean = false
        lateinit var onStatus: OnStatus

        const val notifyID = 1
        const val notifyChannelID = "com.penndev.socks5.vpnService"
        const val notifyChannelName = "Socks5VpnService"
    }


    // tun 设备
    protected var tun: ParcelFileDescriptor? = null

    // 工作进程
    protected var job: Job? = null

    //远程服务器认证
    protected var serviceHost: String = ""
    protected var servicePort: Int = 0
    protected var serviceUser: String = ""
    protected var servicePass: String = ""
    protected var tunDNS: String = "8.8.8.8"
    protected var tunMtu: Int = 1400

    // 启动
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        try {
            val err = setupCommand(intent)
            if (err != null) {
                Toast.makeText(this, err, Toast.LENGTH_SHORT).show()
            } else {
                setupVpnServe()
                onStatus.start()
            }
        } catch (e: Socks5ServiceCloseException) {
            onDestroy()
            Toast.makeText(this, e.message, Toast.LENGTH_SHORT).show()
        } catch (e: Exception) { //处理抛出异常问题
            onDestroy()
            Log.e("penndev", "启动异常", e)
        }
        return START_NOT_STICKY
    }

    // 启动初始化参数
    private fun setupCommand(intent: Intent?): String? {
        if (intent == null) {
            return "传参错误"
        }
        if (intent.getBooleanExtra("close", false)) {
            throw Socks5ServiceCloseException("Closed")
        }
        if (status) {
            return "正在运行中"
        }

        serviceHost = intent.getStringExtra("host")!!

        try {
            InetAddress.getByName(serviceHost)
        } catch (e: Exception) {
            return "错误的host"
        }

        servicePort = intent.getIntExtra("port", 0)
        if (servicePort < 1 || servicePort > 65025) {
            return "错误的port"
        }
        serviceUser = intent.getStringExtra("user")!!
        servicePass = intent.getStringExtra("pass")!!
        return null
    }

    // 启动VPN
    fun setupVpnServe() {
        setupNotifyForeground() //启动通知

        tun = Builder().setMtu(tunMtu).addDnsServer(tunDNS)
            .addRoute("0.0.0.0", 0)
            .addAddress("192.168.1.1", 32)
            .addDisallowedApplication(packageName).establish()

        if (tun == null) {
            throw Socks5ServiceCloseException("获取tun设备失败")
        }
        val tunFd = tun!!.fd.toLong()

        job = GlobalScope.launch {
            try {
                //key.device = "fd://" + // <--- here
                val stack = mobileStack.Stack()
                stack.mtu = tunMtu.toLong()
                stack.srvHost = serviceHost
                stack.srvPort = servicePort.toLong()
                stack.user = serviceUser
                stack.pass = servicePass
                stack.run()
            } catch (e: Exception) { //处理抛出异常问题
                Log.e("penndev", "服务引起异常", e)
            }
        }
    }

    override fun onDestroy() {
        job?.cancel()
        tun?.close()
        onStatus.close()
        stopForeground(true)
        super.onDestroy()
    }

    override fun onRevoke() {
        onDestroy()
    }

    //通知
    protected fun setupNotifyForeground() {
        val resultPendingIntent = PendingIntent.getActivity(
            this, 0,
            Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_NO_CREATE
        )

        var notificationBuilder = NotificationCompat.Builder(this, notifyChannelID)
            .setSmallIcon(R.mipmap.ic_launcher)
            .setContentTitle(getString(R.string.notify_tile_text))
            .setContentText(getString(R.string.notify_context_text))
            .setContentIntent(resultPendingIntent)
            .setPriority(NotificationCompat.PRIORITY_DEFAULT)
        startForeground(notifyID, notificationBuilder.build())
    }
}

