package com.penndev.socks5

import android.app.Activity
import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.content.Intent
import android.net.VpnService
import android.os.Build
import android.os.Bundle
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import com.penndev.socks5.databinding.ActivityMainBinding
import com.penndev.socks5.service.Socks5Service
import com.penndev.socks5.ui.NodeFragment
import com.penndev.socks5.ui.SettingFragment


class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var bindingNode: NodeFragment

    private val activityResultLauncher =
        registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            if (result.resultCode == Activity.RESULT_OK) {
                onStartSocks5Service()
            } else {
                Toast.makeText(this, R.string.toast_main_reject, Toast.LENGTH_LONG).show()
            }
        }

    override fun onCreate(savedInstanceState: Bundle?) {

        binding = ActivityMainBinding.inflate(layoutInflater)
        bindingNode = NodeFragment()

        super.onCreate(savedInstanceState)
        setContentView(binding.root)

        initView()
    }

    private fun initView() {
        fun createNotificationChannel() {
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

        fun onStartUI() {
            binding.mainAction.setImageResource(R.drawable.start)
        }

        fun onStopUI() {
            binding.mainAction.setImageResource(R.drawable.stop)
        }

        fun onStatus(): Socks5Service.OnStatus {
            // 初始化UI状态
            if (Socks5Service.status) onStartUI() else onStopUI()
            return object : Socks5Service.OnStatus {
                override fun start() {
                    super.start()
                    onStartUI()
                }

                override fun stop() {
                    super.stop()
                    onStopUI()
                }
            }
        }

        supportFragmentManager
            .beginTransaction()
            .replace(binding.activitySettingBar.id, SettingFragment())
            .replace(binding.activityNodeBar.id, bindingNode)
            .commit()
        binding.mainAction.setOnClickListener {
            if (Socks5Service.status) onStopSocks5Service() else onStartSocks5Service()
        }
        Socks5Service.onStatus = onStatus()
        createNotificationChannel()
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

        val bundle = Bundle().apply {
            putString("host", bindingNode.Host)
            putInt("port", bindingNode.Port!!.toInt())
            putString("user", bindingNode.user)
            putString("pass", bindingNode.pass)
        }
        startService(Intent(this, Socks5Service::class.java).putExtras(bundle))
    }
}
