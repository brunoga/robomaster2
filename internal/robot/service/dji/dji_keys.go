package dji

import (
	"fmt"
	"reflect"
)

const (
	DJIKeyNone = iota
	DJIProductTest
	DJIProductType
	DJICameraConnection
	DJICameraFirmwareVersion
	DJICameraStartShootPhoto
	DJICameraIsShootingPhoto
	DJICameraPhotoSize
	DJICameraStartRecordVideo
	DJICameraStopRecordVideo
	DJICameraIsRecording
	DJICameraCurrentRecordingTimeInSeconds
	DJICameraVideoFormat
	DJICameraMode
	DJICameraDigitalZoomFactor
	DJICameraAntiFlicker
	DJICameraSwitch
	DJICameraCurrentCameraIndex
	DJICameraHasMainCamera
	DJICameraHasSecondaryCamera
	DJICameraFormatSDCard
	DJICameraSDCardIsFormatting
	DJICameraSDCardIsFull
	DJICameraSDCardHasError
	DJICameraSDCardIsInserted
	DJICameraSDCardTotalSpaceInMB
	DJICameraSDCardRemaingSpaceInMB
	DJICameraSDCardAvailablePhotoCount
	DJICameraSDCardAvailableRecordingTimeInSeconds
	DJICameraIsTimeSynced
	DJICameraDate
	DJICameraVideoTransRate
	DJICameraRequestIFrame
	DJICameraAntiLarsenAlgorithmEnable
	DJIMainControllerConnection
	DJIMainControllerFirmwareVersion
	DJIMainControllerLoaderVersion
	DJIMainControllerVirtualStick
	DJIMainControllerVirtualStickEnabled
	DJIMainControllerChassisSpeedMode
	DJIMainControllerChassisFollowMode
	DJIMainControllerChassisCarControlMode
	DJIMainControllerRecordState
	DJIMainControllerGetRecordSetting
	DJIMainControllerSetRecordSetting
	DJIMainControllerPlayRecordAttr
	DJIMainControllerGetPlayRecordSetting
	DJIMainControllerSetPlayRecordSetting
	DJIMainControllerMaxSpeedForward
	DJIMainControllerMaxSpeedBackward
	DJIMainControllerMaxSpeedLateral
	DJIMainControllerSlopeY
	DJIMainControllerSlopeX
	DJIMainControllerSlopeBreakY
	DJIMainControllerSlopeBreakX
	DJIMainControllerMaxSpeedForwardConfig
	DJIMainControllerMaxSpeedBackwardConfig
	DJIMainControllerMaxSpeedLateralConfig
	DJIMainControllerSlopSpeedYConfig
	DJIMainControllerSlopSpeedXConfig
	DJIMainControllerSlopBreakYConfig
	DJIMainControllerSlopBreakXConfig
	DJIMainControllerChassisPosition
	DJIMainControllerWheelSpeed
	DJIRobomasterMainControllerEscEncodingStatus
	DJIRobomasterMainControllerEscEncodeFlag
	DJIRobomasterMainControllerStartIMUCalibration
	DJIRobomasterMainControllerIMUCalibrationState
	DJIRobomasterMainControllerIMUCalibrationCurrSide
	DJIRobomasterMainControllerIMUCalibrationProgress
	DJIRobomasterMainControllerIMUCalibrationFailCode
	DJIRobomasterMainControllerIMUCalibrationFinishFlag
	DJIRobomasterMainControllerStopIMUCalibration
	DJIRobomasterChassisMode
	DJIRobomasterChassisSpeed
	DJIRobomasterOpenChassisSpeedUpdates
	DJIRobomasterCloseChassisSpeedUpdates
	DJIRobomasterMainControllerRelativePosition
	DJIMainControllerArmServoID
	DJIMainControllerServoAddressing
	DJIRemoteControllerConnection
	DJIGimbalConnection
	DJIGimbalESCFirmwareVersion
	DJIGimbalFirmwareVersion
	DJIGimbalWorkMode
	DJIGimbalControlMode
	DJIGimbalResetPosition
	DJIGimbalResetPositionState
	DJIGimbalCalibration
	DJIGimbalSpeedRotation
	DJIGimbalSpeedRotationEnabled
	DJIGimbalAngleIncrementRotation
	DJIGimbalAngleFrontYawRotation
	DJIGimbalAngleFrontPitchRotation
	DJIGimbalAttitude
	DJIGimbalAutoCalibrate
	DJIGimbalCalibrationStatus
	DJIGimbalCalibrationProgress
	DJIGimbalOpenAttitudeUpdates
	DJIGimbalCloseAttitudeUpdates
	DJIRobomasterSystemConnection
	DJIRobomasterSystemFirmwareVersion
	DJIRobomasterSystemCANFirmwareVersion
	DJIRobomasterSystemScratchFirmwareVersion
	DJIRobomasterSystemSerialNumber
	DJIRobomasterSystemAbilitiesAttack
	DJIRobomasterSystemUnderAbilitiesAttack
	DJIRobomasterSystemKill
	DJIRobomasterSystemRevive
	DJIRobomasterSystemGet1860LinkAck
	DJIMainControllerGetLinkAck
	DJIGimbalGetLinkAck
	DJIRobomasterSystemGameRoleConfig
	DJIRobomasterSystemGameColorConfig
	DJIRobomasterSystemGameStart
	DJIRobomasterSystemGameEnd
	DJIRobomasterSystemDebugLog
	DJIRobomasterSystemSoundEnabled
	DJIRobomasterSystemLeftHeadlightBrightness
	DJIRobomasterSystemRightHeadlightBrightness
	DJIRobomasterSystemLEDColor
	DJIRobomasterSystemUploadScratch
	DJIRobomasterSystemUploadScratchByFTP
	DJIRobomasterSystemUninstallScratchSkill
	DJIRobomasterSystemInstallScratchSkill
	DJIRobomasterSystemInquiryDspMd5
	DJIRobomasterSystemInquiryDspMd5Ack
	DJIRobomasterSystemInquiryDspResourceMd5
	DJIRobomasterSystemInquiryDspResourceMd5Ack
	DJIRobomasterSystemLaunchSinglePlayerCustomSkill
	DJIRobomasterSystemStopSinglePlayerCustomSkill
	DJIRobomasterSystemControlScratch
	DJIRobomasterSystemScratchState
	DJIRobomasterSystemScratchCallback
	DJIRobomasterSystemForesightPosition
	DJIRobomasterSystemPullLogFiles
	DJIRobomasterSystemCurrentHP
	DJIRobomasterSystemTotalHP
	DJIRobomasterSystemCurrentBullets
	DJIRobomasterSystemTotalBullets
	DJIRobomasterSystemEquipments
	DJIRobomasterSystemBuffs
	DJIRobomasterSystemSkillStatus
	DJIRobomasterSystemGunCoolDown
	DJIRobomasterSystemGameConfigList
	DJIRobomasterSystemCarAndSkillID
	DJIRobomasterSystemAppStatus
	DJIRobomasterSystemLaunchMultiPlayerSkill
	DJIRobomasterSystemStopMultiPlayerSkill
	DJIRobomasterSystemConfigSkillTable
	DJIRobomasterSystemWorkingDevices
	DJIRobomasterSystemExceptions
	DJIRobomasterSystemTaskStatus
	DJIRobomasterSystemReturnEnabled
	DJIRobomasterSystemSafeMode
	DJIRobomasterSystemScratchExecuteState
	DJIRobomasterSystemAttitudeInfo
	DJIRobomasterSystemSightBeadPosition
	DJIRobomasterSystemSpeakerLanguage
	DJIRobomasterSystemSpeakerVolumn
	DJIRobomasterSystemChassisSpeedLevel
	DJIRobomasterSystemIsEncryptedFirmware
	DJIRobomasterSystemScratchErrorInfo
	DJIRobomasterSystemScratchOutputInfo
	DJIRobomasterSystemBarrelCoolDown
	DJIRobomasterSystemResetBarrelOverheat
	DJIRobomasterSystemMobileAccelerInfo
	DJIRobomasterSystemMobileGyroAttitudeAngleInfo
	DJIRobomasterSystemMobileGyroRotationRateInfo
	DJIRobomasterSystemEnableAcceleratorSubscribe
	DJIRobomasterSystemEnableGyroRotationRateSubscribe
	DJIRobomasterSystemEnableGyroAttitudeAngleSubscribe
	DJIRobomasterSystemDeactivate
	DJIRobomasterSystemFunctionEnable
	DJIRobomasterSystemIsGameRunning
	DJIRobomasterSystemIsActivated
	DJIRobomasterSystemLowPowerConsumption
	DJIRobomasterSystemEnterLowPowerConsumption
	DJIRobomasterSystemExitLowPowerConsumption
	DJIRobomasterSystemIsLowPowerConsumption
	DJIRobomasterSystemPushFile
	DJIRobomasterSystemPlaySound
	DJIRobomasterSystemPlaySoundStatus
	DJIRobomasterSystemCustomUIAttribute
	DJIRobomasterSystemCustomUIFunctionEvent
	DJIRobomasterSystemTotalMileage
	DJIRobomasterSystemTotalDrivingTime
	DJIRobomasterSystemSetPlayMode
	DJIRobomasterSystemCustomSkillInfo
	DJIRobomasterSystemAddressing
	DJIRobomasterSystemLEDLightEffect
	DJIRobomasterSystemOpenImageTransmission
	DJIRobomasterSystemCloseImageTransmission
	DJIVisionFirmwareVersion
	DJIVisionTrackingAutoLockTarget
	DJIVisionARParameters
	DJIVisionARTagEnabled
	DJIVisionDebugRect
	DJIVisionLaserPosition
	DJIVisionDetectionEnable
	DJIVisionMarkerRunningStatus
	DJIVisionTrackingRunningStatus
	DJIVisionAimbotRunningStatus
	DJIVisionHeadAndShoulderStatus
	DJIVisionHumanDetectionRunningStatus
	DJIVisionUserConfirm
	DJIVisionUserCancel
	DJIVisionUserTrackingRect
	DJIVisionTrackingDistance
	DJIVisionLineColor
	DJIVisionMarkerColor
	DJIVisionMarkerAdvanceStatus
	DJIPerceptionFirmwareVersion
	DJIPerceptionMarkerEnable
	DJIPerceptionMarkerResult
	DJIESCFirmwareVersion1
	DJIESCFirmwareVersion2
	DJIESCFirmwareVersion3
	DJIESCFirmwareVersion4
	DJIESCMotorInfomation1
	DJIESCMotorInfomation2
	DJIESCMotorInfomation3
	DJIESCMotorInfomation4
	DJIWiFiLinkFirmwareVersion
	DJIWiFiLinkDebugInfo
	DJIWiFiLinkMode
	DJIWiFiLinkSSID
	DJIWiFiLinkPassword
	DJIWiFiLinkAvailableChannelNumbers
	DJIWiFiLinkCurrentChannelNumber
	DJIWiFiLinkSNR
	DJIWiFiLinkSNRPushEnabled
	DJIWiFiLinkReboot
	DJIWiFiLinkChannelSelectionMode
	DJIWiFiLinkInterference
	DJIWiFiLinkDeleteNetworkConfig
	DJISDRLinkSNR
	DJISDRLinkBandwidth
	DJISDRLinkChannelSelectionMode
	DJISDRLinkCurrentFreqPoint
	DJISDRLinkCurrentFreqBand
	DJISDRLinkIsDualFreqSupported
	DJISDRLinkUpdateConfigs
	DJIAirLinkConnection
	DJIAirLinkSignalQuality
	DJIAirLinkCountryCode
	DJIAirLinkCountryCodeUpdated
	DJIArmorFirmwareVersion1
	DJIArmorFirmwareVersion2
	DJIArmorFirmwareVersion3
	DJIArmorFirmwareVersion4
	DJIArmorFirmwareVersion5
	DJIArmorFirmwareVersion6
	DJIArmorUnderAttack
	DJIArmorEnterResetID
	DJIArmorCancelResetID
	DJIArmorSkipCurrentID
	DJIArmorResetStatus
	DJIRobomasterWaterGunFirmwareVersion
	DJIRobomasterWaterGunWaterGunFire
	DJIRobomasterWaterGunWaterGunFireWithTimes
	DJIRobomasterWaterGunShootSpeed
	DJIRobomasterWaterGunShootFrequency
	DJIRobomasterInfraredGunConnection
	DJIRobomasterInfraredGunFirmwareVersion
	DJIRobomasterInfraredGunInfraredGunFire
	DJIRobomasterInfraredGunShootFrequency
	DJIRobomasterBatteryFirmwareVersion
	DJIRobomasterBatteryPowerPercent
	DJIRobomasterBatteryVoltage
	DJIRobomasterBatteryTemperature
	DJIRobomasterBatteryCurrent
	DJIRobomasterBatteryShutdown
	DJIRobomasterBatteryReboot
	DJIRobomasterGamePadConnection
	DJIRobomasterGamePadFirmwareVersion
	DJIRobomasterGamePadHasMouse
	DJIRobomasterGamePadHasKeyboard
	DJIRobomasterGamePadCtrlSensitivityX
	DJIRobomasterGamePadCtrlSensitivityY
	DJIRobomasterGamePadCtrlSensitivityYaw
	DJIRobomasterGamePadCtrlSensitivityYawSlop
	DJIRobomasterGamePadCtrlSensitivityYawDeadZone
	DJIRobomasterGamePadCtrlSensitivityPitch
	DJIRobomasterGamePadCtrlSensitivityPitchSlop
	DJIRobomasterGamePadCtrlSensitivityPitchDeadZone
	DJIRobomasterGamePadMouseLeftButton
	DJIRobomasterGamePadMouseRightButton
	DJIRobomasterGamePadC1
	DJIRobomasterGamePadC2
	DJIRobomasterGamePadFire
	DJIRobomasterGamePadFn
	DJIRobomasterGamePadNoCalibrate
	DJIRobomasterGamePadNotAtMiddle
	DJIRobomasterGamePadBatteryWarning
	DJIRobomasterGamePadBatteryPercent
	DJIRobomasterGamePadActivationSettings
	DJIRobomasterGamePadControlEnabled
	DJIRobomasterClawConnection
	DJIRobomasterClawFirmwareVersion
	DJIRobomasterClawCtrl
	DJIRobomasterClawStatus
	DJIRobomasterClawInfoSubscribe
	DJIRobomasterEnableClawInfoSubscribe
	DJIRobomasterArmConnection
	DJIRobomasterArmCtrl
	DJIRobomasterArmCtrlMode
	DJIRobomasterArmCalibration
	DJIRobomasterArmBlockedFlag
	DJIRobomasterArmPositionSubscribe
	DJIRobomasterArmReachLimitX
	DJIRobomasterArmReachLimitY
	DJIRobomasterEnableArmInfoSubscribe
	DJIRobomasterArmControlMode
	DJIRobomasterTOFConnection
	DJIRobomasterTOFLEDColor
	DJIRobomasterTOFOnlineModules
	DJIRobomasterTOFInfoSubscribe
	DJIRobomasterEnableTOFInfoSubscribe
	DJIRobomasterTOFFirmwareVersion1
	DJIRobomasterTOFFirmwareVersion2
	DJIRobomasterTOFFirmwareVersion3
	DJIRobomasterTOFFirmwareVersion4
	DJIRobomasterServoConnection
	DJIRobomasterServoLEDColor
	DJIRobomasterServoSpeed
	DJIRobomasterServoOnlineModules
	DJIRobomasterServoInfoSubscribe
	DJIRobomasterEnableServoInfoSubscribe
	DJIRobomasterServoFirmwareVersion1
	DJIRobomasterServoFirmwareVersion2
	DJIRobomasterServoFirmwareVersion3
	DJIRobomasterServoFirmwareVersion4
	DJIRobomasterSensorAdapterConnection
	DJIRobomasterSensorAdapterOnlineModules
	DJIRobomasterSensorAdapterInfoSubscribe
	DJIRobomasterEnableSensorAdapterInfoSubscribe
	DJIRobomasterSensorAdapterFirmwareVersion1
	DJIRobomasterSensorAdapterFirmwareVersion2
	DJIRobomasterSensorAdapterFirmwareVersion3
	DJIRobomasterSensorAdapterFirmwareVersion4
	DJIRobomasterSensorAdapterFirmwareVersion5
	DJIRobomasterSensorAdapterFirmwareVersion6
	DJIRobomasterSensorAdapterLEDColor
	DJIKeysCount
)

var (
	keyAttributeMap = map[DJIKeys]keyAttributes{
		DJIAirLinkConnection: {117440513, typeof[DJIBoolParamValue](), AccessType_Read},
		// TODO(bga): Add any keys we need here.
	}

	keyByValueMap = map[int]DJIKeys{
		117440513: DJIAirLinkConnection,
		// TODO(bga): Add any keys we 117440513need here.
	}
)

type DJIKeys int

func getKeyAttributeOrPanic(k DJIKeys) keyAttributes {
	ka, ok := keyAttributeMap[k]
	if !ok {
		panic(fmt.Sprintf("Can't get key attributes for key %d.", k))
	}

	return ka
}

func (k DJIKeys) Value() uint32 {
	return getKeyAttributeOrPanic(k).value
}

func (k DJIKeys) DataType() reflect.Type {
	return getKeyAttributeOrPanic(k).dataType
}

func (k DJIKeys) AccessType() AccessType {
	return getKeyAttributeOrPanic(k).accessType
}

type AccessType int

const (
	AccessType_None AccessType = 0
	AccessType_Read AccessType = 1 << (iota - 1)
	AccessType_Write
	AccessType_Action
)

type keyAttributes struct {
	value      uint32
	dataType   reflect.Type
	accessType AccessType
}

func keyByValue(value int) DJIKeys {
	key, ok := keyByValueMap[value]
	if !ok {
		panic(fmt.Sprintf("Can't get key for value %d.", value))
	}

	return key
}

func typeof[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}
