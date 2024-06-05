package com.penndev.socks5.ui.home

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.content.SharedPreferences
import android.net.VpnService
import android.os.Bundle
import android.util.TypedValue
import android.view.Gravity
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import com.google.android.material.snackbar.Snackbar
import com.penndev.socks5.R
import com.penndev.socks5.databinding.FragmentHomeBinding
import com.penndev.socks5.service.Socks5Service

class HomeFragment : Fragment() {
    private var sharedPreferences: SharedPreferences? = null
        get() {
            return context?.getSharedPreferences(context?.packageName, Context.MODE_PRIVATE)
        }

    private lateinit var binding: FragmentHomeBinding
    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        binding = FragmentHomeBinding.inflate(inflater, container, false)
        initView()
        return binding.root
    }

    private fun initView() {
        // 设置节点服务器
        binding.homeNode.setOnClickListener {
            val intent = Intent(activity, HomeNodeActivity::class.java).apply {
                putExtra("close", true)
            }
            startActivity(intent)
        }

        // 处理启动按钮
        fun onStartUI() { binding.mainAction.setImageResource(R.drawable.home_start) }
        fun onStopUI() { binding.mainAction.setImageResource(R.drawable.home_stop) }
        if (Socks5Service.status) onStartUI() else onStopUI()
        Socks5Service.onStatus = object : Socks5Service.OnStatus {
            override fun start() {
                super.start()
                onStartUI()
            }
            override fun stop() {
                super.stop()
                onStopUI()
            }
        }

        var host = sharedPreferences?.getString("host",null)
        if(host != null){
            val port = sharedPreferences?.getInt("port", 1080)
            var user = sharedPreferences?.getString("user", null)
            var pass = sharedPreferences?.getString("pass", null)
            binding.mainAction.setOnClickListener {
                if (Socks5Service.status) onStopSocks5Service() else onStartSocks5Service()
            }
        }else{
            binding.mainAction.setOnClickListener {
                var snackbar = Snackbar.make(binding.root.rootView, R.string.home_srv_no_node, Snackbar.LENGTH_LONG)
                snackbar.show()
            }
        }



    }

    private fun onStopSocks5Service() {
        val intent = Intent(activity, Socks5Service::class.java).apply {
            putExtra("close", true)
        }
        activity?.startService(intent)
    }

    // 授权回调
    private val activityResultLauncher =
        registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            if (result.resultCode == Activity.RESULT_OK) {
                onStartSocks5Service()
            } else {
                Toast.makeText(activity, R.string.permission_vpn_reject, Toast.LENGTH_LONG).show()
            }
        }

    // 启动vpn
    private fun onStartSocks5Service() {
        val intentPrepare = VpnService.prepare(activity)
        if (intentPrepare != null) {
            activityResultLauncher.launch(intentPrepare)
            return
        }

        val bundle = Bundle().apply {
            //putString("host", bindingNode.Host)
            //putInt("port", bindingNode.Port!!.toInt())
            //putString("user", bindingNode.user)
            //putString("pass", bindingNode.pass)
        }


        activity?.startService(Intent(activity, Socks5Service::class.java).putExtras(bundle))
    }

}