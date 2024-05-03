package com.penndev.socks5

import android.app.Activity
import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.content.Intent
import android.content.SharedPreferences
import android.graphics.Color
import android.net.VpnService
import android.os.Build
import android.os.Bundle
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.edit
import com.penndev.socks5.databinding.ActivityMainBackupBinding
import com.penndev.socks5.databinding.ActivityMainBinding
import com.penndev.socks5.service.Socks5Service
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.InputStream
import java.io.OutputStream
import java.net.InetSocketAddress
import java.net.Socket

class MainBackupActivity : AppCompatActivity() {
    // xml UI 实例
    private lateinit var binding: ActivityMainBackupBinding

    // 表单数据持久化
    private lateinit var sharedPreferences: SharedPreferences

    private val activityResultLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
        if (result.resultCode == Activity.RESULT_OK) {
            onStartSocks5Service()
        }else{
            Toast.makeText(this, R.string.toast_main_reject, Toast.LENGTH_LONG).show()
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBackupBinding.inflate(layoutInflater)
        setContentView(binding.root)
        //设置各种初始化
        sharedPreferences = getSharedPreferences(packageName, Context.MODE_PRIVATE)
        binding.inputHost.setText(sharedPreferences.getString("inputHost", ""))
        binding.inputPort.setText(sharedPreferences.getString("inputPortStr", ""))
        binding.inputUser.setText(sharedPreferences.getString("inputUser", ""))
        binding.inputPass.setText(sharedPreferences.getString("inputPass", ""))
        binding.handleActionIcon.setOnClickListener {
            if (Socks5Service.status) onStopSocks5Service() else onStartSocks5Service()
        }
        binding.handleLoadInfo.setOnClickListener {
            CoroutineScope(Dispatchers.IO).launch {
                val message = onCheckSocks5RTT()
                withContext(Dispatchers.Main) {
                    binding.infoLoad.text = message
                }
            }
        }
        Socks5Service.onStatus = onStatus()
        createNotificationChannel()
    }

    // 绑定UI和Socks5Service状态
    private fun onStatus(): Socks5Service.OnStatus {
        // 初始化UI状态
        fun startUI() {
            binding.handleActionIcon.setColorFilter(Color.parseColor("#000000"))
            binding.handleActionText.setText(R.string.activity_main_start_close_text)
        }

        fun closeUI() {
            binding.handleActionIcon.setColorFilter(Color.parseColor("#888888"))
            binding.handleActionText.setText(R.string.activity_main_start_text)
        }
        if (Socks5Service.status) startUI() else closeUI()
        return object : Socks5Service.OnStatus {
            override fun start() {
                super.start()
                startUI()
                // 保存表单数据
                sharedPreferences.edit {
                    putString("inputHost", binding.inputHost.text.toString())
                    putString("inputPortStr", binding.inputPort.text.toString())
                    putString("inputUser", binding.inputUser.text.toString())
                    putString("inputPass", binding.inputPass.text.toString())
                    apply()
                }
            }

            override fun close() {
                super.close()
                closeUI()
            }
        }
    }

    // 尽快创建通知通道。
    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val name = Socks5Service.notifyChannelName
            val descriptionText = getString(R.string.notify_channel_description)
            val importance = NotificationManager.IMPORTANCE_DEFAULT
            val channel = NotificationChannel(Socks5Service.notifyChannelID, name, importance).apply {
                description = descriptionText
            }
            // Register the channel with the system.
            val notificationManager: NotificationManager =
                getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
            notificationManager.createNotificationChannel(channel)
        }
    }

    private fun onStopSocks5Service() {
        startService(Intent(this, Socks5Service::class.java).putExtra("close", true))
    }

    // 启动 Socks5Service的详细步骤和流程。
    private fun onStartSocks5Service() {
        val intentPrepare = VpnService.prepare(this)
        if (intentPrepare != null) {
            activityResultLauncher.launch(intentPrepare)
            return
        }

        val inputHost = binding.inputHost.text.toString()
        val inputPortStr = binding.inputPort.text.toString()
        val inputUser = binding.inputUser.text.toString()
        val inputPass = binding.inputPass.text.toString()
        if (inputHost == "" || inputPortStr == "") {
            Toast.makeText(this, getString(R.string.toast_main_param_error), Toast.LENGTH_LONG).show()
            return
        }
        val inputPort = inputPortStr.toInt()

        val bundle = Bundle().apply {
            putString("host", inputHost)
            putInt("port", inputPort)
            putString("user", inputUser)
            putString("pass", inputPass)
        }
        startService(Intent(this, Socks5Service::class.java).putExtras(bundle))
    }

    // socks5 探测登录方式
    private fun onCheckSocks5RTT(): String {
        var msg: String
        try {
            msg = Socket().use { socket ->
                val host: String = binding.inputHost.text.toString()
                val port: Int = binding.inputPort.text.toString().toInt()
                val timeout = 10000
                socket.connect(InetSocketAddress(host, port), timeout)
                socket.soTimeout = timeout
                val outputStream: OutputStream = socket.getOutputStream()
                val inputStream: InputStream = socket.getInputStream()

                val startTime = System.currentTimeMillis()
                outputStream.write(byteArrayOf(0x05, 0x01, 0x00)); outputStream.flush() // 无密码探测

                val packet = ByteArray(2)
                val bytesRead = inputStream.read(packet)
                val endTime = System.currentTimeMillis()

                if (bytesRead == 2 && packet[0].toInt() == 0x05) {
                    getString(R.string.activity_main_try_srv_succeed, endTime - startTime)
                } else {
                    getString(R.string.activity_main_try_srv_fail, packet.joinToString(" |B"))
                }
            }
        } catch (e: Exception) {
            msg =  getString(R.string.activity_main_try_srv_fail, e.message)
        }
        return msg
    }
}
