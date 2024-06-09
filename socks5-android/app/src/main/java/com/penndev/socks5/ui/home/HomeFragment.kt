package com.penndev.socks5.ui.home

import android.app.Activity
import android.content.Intent
import android.net.VpnService
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.Fragment
import com.penndev.socks5.R
import com.penndev.socks5.databinding.FragmentHomeBinding
import com.penndev.socks5.databinding.NodeData
import com.penndev.socks5.service.Socks5Service
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.net.InetSocketAddress
import java.net.Socket

class HomeFragment : Fragment() {

    private lateinit var binding: FragmentHomeBinding
    private lateinit var nodedata: NodeData

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        binding = FragmentHomeBinding.inflate(inflater, container, false)
        nodedata = NodeData(requireContext())
        createView()
        return binding.root
    }

    fun onStartUI() { binding.mainAction.setImageResource(R.drawable.home_start) }

    fun onStopUI() { binding.mainAction.setImageResource(R.drawable.home_stop) }

    private fun createView() {
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
        binding.homeNode.setOnClickListener {
            startActivity(Intent(activity, HomeNodeActivity::class.java))}
        binding.mainAction.setOnClickListener {
            if (Socks5Service.status) onStopSocks5Service() else onStartSocks5Service()
        }
        initView()
    }

    override fun onStart(){
        super.onStart()
        initView()
    }

    private fun initView() {
        if (Socks5Service.status) onStartUI() else onStopUI()
        var host = nodedata.host
        var port = nodedata.port
        if(host != null && port!! > 0){
            val hoststr = "${nodedata.host}:${nodedata.port}"
            binding.currentNodeHost.setText(hoststr)
            CoroutineScope(Dispatchers.IO).launch {
                var ms = hostSocks5RTT(host!!, port!!.toInt())
                withContext(Dispatchers.Main) {
                    binding.currentNodeSpeed.text = "${ms}ms"
                }
            }
        }

    }

    private fun hostSocks5RTT(host:String,port:Int): Int {
        try {
            return Socket().use { socket ->
                val timeout = 30000
                socket.connect(InetSocketAddress(host, port), timeout)
                socket.soTimeout = timeout
                val outputStream = socket.getOutputStream()
                val inputStream = socket.getInputStream()

                val startTime = System.currentTimeMillis()
                outputStream.write(byteArrayOf(0x05, 0x01, 0x00)); outputStream.flush() // 无密码探测

                val packet = ByteArray(2)
                val bytesRead = inputStream.read(packet)
                val endTime = System.currentTimeMillis()

                if (bytesRead == 2 && packet[0].toInt() == 0x05) {
                    return (endTime - startTime).toInt()
                } else {
                    return 0
                }
            }
        } catch (e: Exception) {
            return 0
        }
        return 0
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
        if(nodedata.host == null || nodedata.host == "") {
            var alert = AlertDialog.Builder(requireContext())
            alert.setMessage(R.string.home_node_current_host)
            alert.show()
            return
        }

        val intentPrepare = VpnService.prepare(activity)
        if (intentPrepare != null) {
            activityResultLauncher.launch(intentPrepare)
            return
        }

        val bundle = Bundle().apply {
            putString("host", nodedata.host)
            putInt("port", nodedata.port!!)
            putString("user", nodedata.user)
            putString("pass", nodedata.pass)
        }
        activity?.startService(Intent(activity, Socks5Service::class.java).putExtras(bundle))
    }

}