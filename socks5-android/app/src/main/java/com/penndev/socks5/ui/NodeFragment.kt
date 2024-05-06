package com.penndev.socks5.ui

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import com.google.android.material.bottomsheet.BottomSheetDialog
import com.penndev.socks5.R
import com.penndev.socks5.databinding.FragmentNodeBinding


class NodeFragment : Fragment() {

    private lateinit var binding: FragmentNodeBinding

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        binding = FragmentNodeBinding.inflate(inflater, container, false)
        initView()
        return binding.root
    }

    private fun initView() {
        // 加载 Fragment 的布局文件
        val bottomSheetDialog = BottomSheetDialog(requireContext())
        val bottomSheetView: View = layoutInflater.inflate(R.layout.dialog_node, null)
        bottomSheetDialog.setContentView(bottomSheetView)
        binding.root.setOnClickListener{
            bottomSheetDialog.show()
        }
    }
}
