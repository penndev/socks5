package com.penndev.socks5.ui

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import com.google.android.material.bottomsheet.BottomSheetDialog
import com.penndev.socks5.R
import com.penndev.socks5.databinding.FragmentSettingBinding


class SettingFragment : Fragment() {

    private lateinit var binding: FragmentSettingBinding

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        binding = FragmentSettingBinding.inflate(inflater, container, false)
        initView()
        return binding.root
    }

    private fun initView() {
        val bottomSheetDialog = BottomSheetDialog(requireContext())
        val bottomSheetView: View = layoutInflater.inflate(R.layout.dialog_setting, null)
        bottomSheetDialog.setContentView(bottomSheetView)
        binding.openDialogSetting.setOnClickListener{
            bottomSheetDialog.show()
        }
    }
}
