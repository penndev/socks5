package com.penndev.socks5.service

import android.annotation.SuppressLint
import android.app.PendingIntent
import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor
import android.util.Log
import android.widget.Toast
import androidx.core.app.NotificationCompat
import com.penndev.socks5.MainActivity
import com.penndev.socks5.R
import com.penndev.socks5.databinding.NodeData
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import java.net.InetAddress

class Socks5ServiceCloseException(message: String) : Exception(message)

class Socks5Service : VpnService() {

    // 回调UI状态
    interface OnStatus {
        fun start() {
            status = true
        }

        fun stop() {
            status = false
        }
    }

    companion object {
        var status: Boolean = false
        lateinit var onStatus: OnStatus

        const val notifyID = 1
        const val notifyChannelID = "com.penndev.socks5.vpnService"
        const val notifyChannelName = "Socks5VpnService"
        var readByteLen: Long = 0
        var writeByteLen: Long = 0
    }

    // tun 设备
    private var tun: ParcelFileDescriptor? = null

    // 工作进程
    private var job: Job? = null

    //远程服务器认证
    //private var serviceHost: String = ""
    //private var servicePort: Int = 0
    //private var serviceUser: String = ""
    //private var servicePass: String = ""
    private lateinit var nodeData: NodeData

    private var tunDNS: String = "8.8.8.8"
    private var tunMtu: Int = 1400

    // 启动
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        try {
            nodeData = NodeData(this)
            val err = setupCommand(intent)
            if (err != null) {
                Toast.makeText(this, err, Toast.LENGTH_SHORT).show()
            } else {
                setupVpnServe()
                onStatus.start()
            }
        } catch (e: Socks5ServiceCloseException) {
            onDestroy()
        } catch (e: Exception) { //处理抛出异常问题
            onDestroy()
            Log.e("penndev", "启动异常", e)
        }
        return START_NOT_STICKY
    }

    // 启动初始化参数
    private fun setupCommand(intent: Intent?): String? {
        if (intent == null) {
            return getString(R.string.toast_service_param_error)
        }
        if (intent.getBooleanExtra("status", false)) {
            throw Socks5ServiceCloseException("Closed")
        }
        if (status) {
            return getString(R.string.toast_service_always_run)
        }

        try {
            InetAddress.getByName(nodeData.host)
        } catch (e: Exception) {
            Log.e("debug", "nodeData.host ${nodeData.host} err", e)
            return getString(R.string.toast_service_host_error)
        }

        val servicePort = nodeData.port!!
        if (servicePort < 1 || servicePort > 65025) {
            return getString(R.string.toast_service_port_error)
        }
        return null
    }

    // 启动VPN
    private fun setupVpnServe() {
        setupNotifyForeground() //启动通知
        val tunDevice = Builder()
            .setMtu(tunMtu)
            .addDnsServer(tunDNS)
            .addRoute("0.0.0.0", 0)
            .addAddress("192.168.0.1", 32)
            .addDisallowedApplication(packageName).establish()
            ?: throw Socks5ServiceCloseException(getString(R.string.toast_service_tun_null))
        tun = tunDevice

        val tunFd = tunDevice.fd.toLong()
        job = GlobalScope.launch {
            try {
                val stack = mobileStack.Stack()
                stack.tunFd = tunFd
                stack.mtu = tunMtu.toLong()
                stack.srvHost = nodeData.host
                stack.srvPort = nodeData.port!!.toLong()
                stack.user = nodeData.user
                stack.pass = nodeData.pass
                if(nodeData.typeTcpEnable == true) {
                    stack.tcpEnable = true
                }
                if(nodeData.typeUdpEnable == true) {
                    stack.udpEnable = true
                }
                stack.handle = object : mobileStack.StackHandle{
                    override fun error(err: String?) {
                        err?.let { Log.e("go_error", it) }
                    }
                    override fun readLen(i: Long) {
                        writeByteLen += i
                    }
                    override fun writeLen(i: Long) {
                        readByteLen += i
                    }
                }
                val status = stack.run()
                if (!status) {
                    delay(200)
                    onDestroy()
                }
            } catch (e: Exception) { //处理抛出异常问题
                Log.e("debug", "mobileStack.Stack", e)
            }
        }
    }

    override fun onDestroy() {
        job?.cancel()
        tun?.close()
        onStatus.stop()
        stopForeground(true)
        super.onDestroy()
    }

    override fun onRevoke() {
        onDestroy()
    }

    //通知
    @SuppressLint("UnspecifiedImmutableFlag", "ForegroundServiceType")
    private fun setupNotifyForeground() {
        val intent = Intent(this, MainActivity::class.java)
        intent.action = Intent.ACTION_MAIN
        intent.addCategory(Intent.CATEGORY_LAUNCHER)
        val resultPendingIntent = PendingIntent.getActivity(
            this, 0, intent, PendingIntent.FLAG_IMMUTABLE
        )

        val notificationBuilder = NotificationCompat.Builder(this, notifyChannelID)
            .setSmallIcon(R.mipmap.ic_launcher)
            .setContentTitle(getString(R.string.notify_tile_text))
            .setContentText(getString(R.string.notify_context_text))
            .setContentIntent(resultPendingIntent)
            .setPriority(NotificationCompat.PRIORITY_DEFAULT)
        startForeground(notifyID, notificationBuilder.build())
    }
}
