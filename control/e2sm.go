/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/

package control

/*
#include <e2sm/wrapper.h>
#cgo LDFLAGS: -le2smwrapper
#cgo CFLAGS: -I/usr/local/include/e2sm
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"unsafe"
	"fmt"
)

type E2sm struct {
}

func (c *E2sm) SetEventTriggerDefinition(buffer []byte, eventTriggerCount int, RTPeriods []int64) (newBuffer []byte, err error) {
	cptr := unsafe.Pointer(&buffer[0])
	periods := unsafe.Pointer(&RTPeriods[0])
	size := C.e2sm_encode_ric_event_trigger_definition(cptr, C.size_t(len(buffer)), C.size_t(eventTriggerCount), (*C.long)(periods))
	if size < 0 {
		return make([]byte, 0), errors.New("e2sm wrapper is unable to set EventTriggerDefinition due to wrong or invalid input")
	}
	newBuffer = C.GoBytes(cptr, (C.int(size)+7)/8)
	return
}

func (c *E2sm) SetActionDefinition(buffer []byte, ricStyleType int64) (newBuffer []byte, err error) {
	cptr := unsafe.Pointer(&buffer[0])
	size := C.e2sm_encode_ric_action_definition(cptr, C.size_t(len(buffer)), C.long(ricStyleType))
	if size < 0 {
		return make([]byte, 0), errors.New("e2sm wrapper is unable to set ActionDefinition due to wrong or invalid input")
	}
	newBuffer = C.GoBytes(cptr, (C.int(size)+7)/8)
	return
}

func (c *E2sm) GetIndicationHeader(buffer []byte) (indHdr *IndicationHeader, err error) {
	fmt.Println("////////////entered e2sm GetIndicationHeader")
	cptr := unsafe.Pointer(&buffer[0])
	//fmt.Println("/////IndicationHeader= %x",IndicationHeader)
	indHdr = &IndicationHeader{}
	fmt.Println("///////////E2sm////indHdr = &IndicationHeader= %x",indHdr)
	decodedHdr := C.e2sm_decode_ric_indication_header(cptr, C.size_t(len(buffer)))
	
	fmt.Println("///////////E2sm////decodedHdr= %x",decodedHdr)
	if decodedHdr == nil {
		return indHdr, errors.New("e2sm wrapper is unable to get IndicationHeader due to wrong or invalid input")
	}
	defer C.e2sm_free_ric_indication_header(decodedHdr)

	indHdr.IndHdrType = int32(decodedHdr.present)
	fmt.Println("////indHdr.IndHdrType= %", indHdr.IndHdrType)
	if indHdr.IndHdrType == 1 {
		fmt.Println("////entered indHdr.IndHdrType= 1")
		indHdrFormat1 := &IndicationHeaderFormat1{}
		indHdrFormat1_C := *(**C.E2SM_KPM_IndicationHeader_Format1_t)(unsafe.Pointer(&decodedHdr.choice[0]))

		if indHdrFormat1_C.id_GlobalKPMnode_ID != nil {
			globalKPMnodeID_C := (*C.GlobalKPMnode_ID_t)(indHdrFormat1_C.id_GlobalKPMnode_ID)

			indHdrFormat1.GlobalKPMnodeIDType = int32(globalKPMnodeID_C.present)
			fmt.Println("////indHdrFormat1.GlobalKPMnodeIDType= %",indHdrFormat1.GlobalKPMnodeIDType)
			if indHdrFormat1.GlobalKPMnodeIDType == 1 {
				globalgNBID := &GlobalKPMnodegNBIDType{}
				globalgNBID_C := (*C.GlobalKPMnode_gNB_ID_t)(unsafe.Pointer(&globalKPMnodeID_C.choice[0]))

				plmnID_C := globalgNBID_C.global_gNB_ID.plmn_id
				globalgNBID.GlobalgNBID.PlmnID.Buf = C.GoBytes(unsafe.Pointer(plmnID_C.buf), C.int(plmnID_C.size))
				globalgNBID.GlobalgNBID.PlmnID.Size = int(plmnID_C.size)

				globalgNBID_gNBID_C := globalgNBID_C.global_gNB_ID.gnb_id
				globalgNBID.GlobalgNBID.GnbIDType = int(globalgNBID_gNBID_C.present)
				if globalgNBID.GlobalgNBID.GnbIDType == 1 {
					gNBID := &GNBID{}
					gNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globalgNBID_gNBID_C.choice[0]))

					gNBID.Buf = C.GoBytes(unsafe.Pointer(gNBID_C.buf), C.int(gNBID_C.size))
					gNBID.Size = int(gNBID_C.size)
					gNBID.BitsUnused = int(gNBID_C.bits_unused)

					globalgNBID.GlobalgNBID.GnbID = gNBID
					fmt.Println("///in type 1 globalgNBID.GlobalgNBID.GnbID = gNBID= %", globalgNBID.GlobalgNBID.GnbID)
				}

				if globalgNBID_C.gNB_CU_UP_ID != nil {
					globalgNBID.GnbCUUPID = &Integer{}
					globalgNBID.GnbCUUPID.Buf = C.GoBytes(unsafe.Pointer(globalgNBID_C.gNB_CU_UP_ID.buf), C.int(globalgNBID_C.gNB_CU_UP_ID.size))
					globalgNBID.GnbCUUPID.Size = int(globalgNBID_C.gNB_CU_UP_ID.size)
				}

				if globalgNBID_C.gNB_DU_ID != nil {
					globalgNBID.GnbDUID = &Integer{}
					globalgNBID.GnbDUID.Buf = C.GoBytes(unsafe.Pointer(globalgNBID_C.gNB_DU_ID.buf), C.int(globalgNBID_C.gNB_DU_ID.size))
					globalgNBID.GnbDUID.Size = int(globalgNBID_C.gNB_DU_ID.size)
				}

				indHdrFormat1.GlobalKPMnodeID = globalgNBID
				fmt.Println("///in type 1 indHdrFormat1.GlobalKPMnodeID = globalgNBID= %", indHdrFormat1.GlobalKPMnodeID)
			} else if indHdrFormat1.GlobalKPMnodeIDType == 2 {
				fmt.Println("////entered else if type 2///")
				globalengNBID := &GlobalKPMnodeengNBIDType{}
				globalengNBID_C := (*C.GlobalKPMnode_en_gNB_ID_t)(unsafe.Pointer(&globalKPMnodeID_C.choice[0]))

				plmnID_C := globalengNBID_C.global_gNB_ID.pLMN_Identity
				globalengNBID.PlmnID.Buf = C.GoBytes(unsafe.Pointer(plmnID_C.buf), C.int(plmnID_C.size))
				globalengNBID.PlmnID.Size = int(plmnID_C.size)

				globalengNBID_gNBID_C := globalengNBID_C.global_gNB_ID.gNB_ID
				globalengNBID.GnbIDType = int(globalengNBID_gNBID_C.present)
				if globalengNBID.GnbIDType == 1 {
					engNBID := &ENGNBID{}
					engNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globalengNBID_gNBID_C.choice[0]))

					engNBID.Buf = C.GoBytes(unsafe.Pointer(engNBID_C.buf), C.int(engNBID_C.size))
					engNBID.Size = int(engNBID_C.size)
					engNBID.BitsUnused = int(engNBID_C.bits_unused)

					globalengNBID.GnbID = engNBID
					fmt.Println("///in type 2 globalengNBID.GnbID = engNBID= %", globalengNBID.GnbID)
				}

				indHdrFormat1.GlobalKPMnodeID = globalengNBID
				fmt.Println("///in type 2 indHdrFormat1.GlobalKPMnodeID = globalengNBID= %", indHdrFormat1.GlobalKPMnodeID )
			} else if indHdrFormat1.GlobalKPMnodeIDType == 3 {
				fmt.Println("////entered else if type 3///")
				globalngeNBID := &GlobalKPMnodengeNBIDType{}
				globalngeNBID_C := (*C.GlobalKPMnode_ng_eNB_ID_t)(unsafe.Pointer(&globalKPMnodeID_C.choice[0]))

				plmnID_C := globalngeNBID_C.global_ng_eNB_ID.plmn_id
				globalngeNBID.PlmnID.Buf = C.GoBytes(unsafe.Pointer(plmnID_C.buf), C.int(plmnID_C.size))
				globalngeNBID.PlmnID.Size = int(plmnID_C.size)

				globalngeNBID_eNBID_C := globalngeNBID_C.global_ng_eNB_ID.enb_id
				globalngeNBID.EnbIDType = int(globalngeNBID_eNBID_C.present)
				if globalngeNBID.EnbIDType == 1 {
					ngeNBID := &NGENBID_Macro{}
					ngeNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globalngeNBID_eNBID_C.choice[0]))

					ngeNBID.Buf = C.GoBytes(unsafe.Pointer(ngeNBID_C.buf), C.int(ngeNBID_C.size))
					ngeNBID.Size = int(ngeNBID_C.size)
					ngeNBID.BitsUnused = int(ngeNBID_C.bits_unused)

					globalngeNBID.EnbID = ngeNBID
				} else if globalngeNBID.EnbIDType == 2 {
					ngeNBID := &NGENBID_ShortMacro{}
					ngeNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globalngeNBID_eNBID_C.choice[0]))

					ngeNBID.Buf = C.GoBytes(unsafe.Pointer(ngeNBID_C.buf), C.int(ngeNBID_C.size))
					ngeNBID.Size = int(ngeNBID_C.size)
					ngeNBID.BitsUnused = int(ngeNBID_C.bits_unused)

					globalngeNBID.EnbID = ngeNBID
				} else if globalngeNBID.EnbIDType == 3 {
					ngeNBID := &NGENBID_LongMacro{}
					ngeNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globalngeNBID_eNBID_C.choice[0]))

					ngeNBID.Buf = C.GoBytes(unsafe.Pointer(ngeNBID_C.buf), C.int(ngeNBID_C.size))
					ngeNBID.Size = int(ngeNBID_C.size)
					ngeNBID.BitsUnused = int(ngeNBID_C.bits_unused)

					globalngeNBID.EnbID = ngeNBID
					fmt.Println("///in type 3 globalngeNBID.EnbID = ngeNBID= %", globalngeNBID.EnbID)
				}

				indHdrFormat1.GlobalKPMnodeID = globalngeNBID
				fmt.Println("///in type 3 indHdrFormat1.GlobalKPMnodeID = globalngeNBID= %", indHdrFormat1.GlobalKPMnodeID)
			} else if indHdrFormat1.GlobalKPMnodeIDType == 4 {
				fmt.Println("////entered else if type 4///")
				globaleNBID := &GlobalKPMnodeeNBIDType{}
				globaleNBID_C := (*C.GlobalKPMnode_eNB_ID_t)(unsafe.Pointer(&globalKPMnodeID_C.choice[0]))

				plmnID_C := globaleNBID_C.global_eNB_ID.pLMN_Identity
				globaleNBID.PlmnID.Buf = C.GoBytes(unsafe.Pointer(plmnID_C.buf), C.int(plmnID_C.size))
				globaleNBID.PlmnID.Size = int(plmnID_C.size)

				globaleNBID_eNBID_C := globaleNBID_C.global_eNB_ID.eNB_ID
				globaleNBID.EnbIDType = int(globaleNBID_eNBID_C.present)
				if globaleNBID.EnbIDType == 1 {
					eNBID := &ENBID_Macro{}
					eNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globaleNBID_eNBID_C.choice[0]))

					eNBID.Buf = C.GoBytes(unsafe.Pointer(eNBID_C.buf), C.int(eNBID_C.size))
					eNBID.Size = int(eNBID_C.size)
					eNBID.BitsUnused = int(eNBID_C.bits_unused)

					globaleNBID.EnbID = eNBID
				} else if globaleNBID.EnbIDType == 2 {
					eNBID := &ENBID_Home{}
					eNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globaleNBID_eNBID_C.choice[0]))

					eNBID.Buf = C.GoBytes(unsafe.Pointer(eNBID_C.buf), C.int(eNBID_C.size))
					eNBID.Size = int(eNBID_C.size)
					eNBID.BitsUnused = int(eNBID_C.bits_unused)

					globaleNBID.EnbID = eNBID
				} else if globaleNBID.EnbIDType == 3 {
					eNBID := &ENBID_ShortMacro{}
					eNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globaleNBID_eNBID_C.choice[0]))

					eNBID.Buf = C.GoBytes(unsafe.Pointer(eNBID_C.buf), C.int(eNBID_C.size))
					eNBID.Size = int(eNBID_C.size)
					eNBID.BitsUnused = int(eNBID_C.bits_unused)

					globaleNBID.EnbID = eNBID
				} else if globaleNBID.EnbIDType == 4 {
					eNBID := &ENBID_LongMacro{}
					eNBID_C := (*C.BIT_STRING_t)(unsafe.Pointer(&globaleNBID_eNBID_C.choice[0]))

					eNBID.Buf = C.GoBytes(unsafe.Pointer(eNBID_C.buf), C.int(eNBID_C.size))
					eNBID.Size = int(eNBID_C.size)
					eNBID.BitsUnused = int(eNBID_C.bits_unused)

					globaleNBID.EnbID = eNBID
					fmt.Println("///in type 4 globaleNBID.EnbID = ngeNBID= %", globaleNBID.EnbID)
				}

				indHdrFormat1.GlobalKPMnodeID = globaleNBID
				fmt.Println("///in type 4 indHdrFormat1.GlobalKPMnodeID = globaleNBID= %", indHdrFormat1.GlobalKPMnodeID)
			}
		} else {
			fmt.Println("/// entered else indHdrFormat1.GlobalKPMnodeIDType = 0")
			indHdrFormat1.GlobalKPMnodeIDType = 0
		}

		if indHdrFormat1_C.nRCGI != nil {
			indHdrFormat1.NRCGI = &NRCGIType{}

			plmnID := indHdrFormat1_C.nRCGI.pLMN_Identity
			indHdrFormat1.NRCGI.PlmnID.Buf = C.GoBytes(unsafe.Pointer(plmnID.buf), C.int(plmnID.size))
			indHdrFormat1.NRCGI.PlmnID.Size = int(plmnID.size)

			nRCellID := indHdrFormat1_C.nRCGI.nRCellIdentity
			indHdrFormat1.NRCGI.NRCellID.Buf = C.GoBytes(unsafe.Pointer(nRCellID.buf), C.int(nRCellID.size))
			indHdrFormat1.NRCGI.NRCellID.Size = int(nRCellID.size)
			indHdrFormat1.NRCGI.NRCellID.BitsUnused = int(nRCellID.bits_unused)
		}

		if indHdrFormat1_C.pLMN_Identity != nil {
			indHdrFormat1.PlmnID = &OctetString{}

			indHdrFormat1.PlmnID.Buf = C.GoBytes(unsafe.Pointer(indHdrFormat1_C.pLMN_Identity.buf), C.int(indHdrFormat1_C.pLMN_Identity.size))
			indHdrFormat1.PlmnID.Size = int(indHdrFormat1_C.pLMN_Identity.size)
		}

		if indHdrFormat1_C.sliceID != nil {
			indHdrFormat1.SliceID = &SliceIDType{}

			sST := indHdrFormat1_C.sliceID.sST
			indHdrFormat1.SliceID.SST.Buf = C.GoBytes(unsafe.Pointer(sST.buf), C.int(sST.size))
			indHdrFormat1.SliceID.SST.Size = int(sST.size)

			if indHdrFormat1_C.sliceID.sD != nil {
				indHdrFormat1.SliceID.SD = &OctetString{}

				sD := indHdrFormat1_C.sliceID.sD
				indHdrFormat1.SliceID.SD.Buf = C.GoBytes(unsafe.Pointer(sD.buf), C.int(sD.size))
				indHdrFormat1.SliceID.SD.Size = int(sD.size)
			}
		}
		fmt.Println("////before if ! nil indHdrFormat1.FiveQI = %", indHdrFormat1.FiveQI)

		if indHdrFormat1_C.fiveQI != nil {
			indHdrFormat1.FiveQI = *(*int64)(unsafe.Pointer(indHdrFormat1_C.fiveQI))
		} else {
			indHdrFormat1.FiveQI = -1
		}
		fmt.Println("////before if ! nil indHdrFormat1_C.qci= %", indHdrFormat1_C.qci)

		if indHdrFormat1_C.qci != nil {
			indHdrFormat1.Qci = *(*int64)(unsafe.Pointer(indHdrFormat1_C.qci))
		} else {
			indHdrFormat1.Qci = -1
		}

		indHdr.IndHdr = indHdrFormat1
		fmt.Println("///at the end before else///indHdr.IndHdr = %", indHdr.IndHdr)
	} else {
		return indHdr, errors.New("Unknown RIC Indication Header type")
	}

	return
}

func (c *E2sm) GetIndicationMessage(buffer []byte) (indMsg *IndicationMessage, err error) {
	fmt.Println("////////////entered e2sm getindicationmessage")
	cptr := unsafe.Pointer(&buffer[0])
	indMsg = &IndicationMessage{}
	fmt.Println("///////e2sm indMsg= IndicationMessage %x", indMsg)
	decodedMsg := C.e2sm_decode_ric_indication_message(cptr, C.size_t(len(buffer)))
	if decodedMsg == nil {
		fmt.Println("/////e2sm decodeMsg == nil")
		return indMsg, errors.New("e2sm wrapper is unable to get IndicationMessage due to wrong or invalid input")
	}
	defer C.e2sm_free_ric_indication_message(decodedMsg)

	indMsg.StyleType = int64(decodedMsg.ric_Style_Type)

	indMsg.IndMsgType = int32(decodedMsg.indicationMessage.present)
	fmt.Println("///////e2sm indMsg.IndMsgType %x", indMsg.IndMsgType)

	if indMsg.IndMsgType == 1 {
		fmt.Println("////////////entered e2sm if indMsg.IndMsgType == 1")
		indMsgFormat1 := &IndicationMessageFormat1{}
		indMsgFormat1_C := *(**C.E2SM_KPM_IndicationMessage_Format1_t)(unsafe.Pointer(&decodedMsg.indicationMessage.choice[0]))

		indMsgFormat1.PMContainerCount = int(indMsgFormat1_C.pm_Containers.list.count)
		for i := 0; i < indMsgFormat1.PMContainerCount; i++ {
			fmt.Println("---------------e2sm 1st for loop i= %d", i)
			pmContainer := &indMsgFormat1.PMContainers[i]
			fmt.Println("////////////e2sm pmContainer= &indMsgFormat1.PMContainers i %x", pmContainer)
			var sizeof_PM_Containers_List_t *C.PM_Containers_List_t
			pmContainer_C := *(**C.PM_Containers_List_t)(unsafe.Pointer(uintptr(unsafe.Pointer(indMsgFormat1_C.pm_Containers.list.array)) + (uintptr)(i)*unsafe.Sizeof(sizeof_PM_Containers_List_t)))
			fmt.Println("////////////e2sm pmContainer= %x", pmContainer_C)

			if pmContainer_C.performanceContainer != nil {
				fmt.Println("////////////entered e2sm if pmContainer_C.performanceContainer != nil")
				pfContainer := &PFContainerType{}

				pfContainer.ContainerType = int32(pmContainer_C.performanceContainer.present)
				fmt.Println("///e2sm pfContainer.ContainerType %d", pfContainer.ContainerType)

				if pfContainer.ContainerType == 1 {
					fmt.Println("////////////entered e2sm if pfContainer.ContainerType == 1")
					oDU_PF := &ODUPFContainerType{}
					oDU_PF_C := *(**C.ODU_PF_Container_t)(unsafe.Pointer(&pmContainer_C.performanceContainer.choice[0]))

					oDU_PF.CellResourceReportCount = int(oDU_PF_C.cellResourceReportList.list.count)
					for j := 0; j < oDU_PF.CellResourceReportCount; j++ {
						cellResourceReport := &oDU_PF.CellResourceReports[j]
						var sizeof_CellResourceReportListItem_t *C.CellResourceReportListItem_t
						cellResourceReport_C := *(**C.CellResourceReportListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(oDU_PF_C.cellResourceReportList.list.array)) + (uintptr)(j)*unsafe.Sizeof(sizeof_CellResourceReportListItem_t)))

						cellResourceReport.NRCGI.PlmnID.Buf = C.GoBytes(unsafe.Pointer(cellResourceReport_C.nRCGI.pLMN_Identity.buf), C.int(cellResourceReport_C.nRCGI.pLMN_Identity.size))
						cellResourceReport.NRCGI.PlmnID.Size = int(cellResourceReport_C.nRCGI.pLMN_Identity.size)

						cellResourceReport.NRCGI.NRCellID.Buf = C.GoBytes(unsafe.Pointer(cellResourceReport_C.nRCGI.nRCellIdentity.buf), C.int(cellResourceReport_C.nRCGI.nRCellIdentity.size))
						cellResourceReport.NRCGI.NRCellID.Size = int(cellResourceReport_C.nRCGI.nRCellIdentity.size)
						cellResourceReport.NRCGI.NRCellID.BitsUnused = int(cellResourceReport_C.nRCGI.nRCellIdentity.bits_unused)

						if cellResourceReport_C.dl_TotalofAvailablePRBs != nil {
							cellResourceReport.TotalofAvailablePRBs.DL = int64(*cellResourceReport_C.dl_TotalofAvailablePRBs)
						} else {
							cellResourceReport.TotalofAvailablePRBs.DL = -1
						}

						if cellResourceReport_C.ul_TotalofAvailablePRBs != nil {
							cellResourceReport.TotalofAvailablePRBs.UL = int64(*cellResourceReport_C.ul_TotalofAvailablePRBs)
						} else {
							cellResourceReport.TotalofAvailablePRBs.UL = -1
						}

						cellResourceReport.ServedPlmnPerCellCount = int(cellResourceReport_C.servedPlmnPerCellList.list.count)
						for k := 0; k < cellResourceReport.ServedPlmnPerCellCount; k++ {
							servedPlmnPerCell := cellResourceReport.ServedPlmnPerCells[k]
							var sizeof_ServedPlmnPerCellListItem_t *C.ServedPlmnPerCellListItem_t
							servedPlmnPerCell_C := *(**C.ServedPlmnPerCellListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(cellResourceReport_C.servedPlmnPerCellList.list.array)) + (uintptr)(k)*unsafe.Sizeof(sizeof_ServedPlmnPerCellListItem_t)))

							servedPlmnPerCell.PlmnID.Buf = C.GoBytes(unsafe.Pointer(servedPlmnPerCell_C.pLMN_Identity.buf), C.int(servedPlmnPerCell_C.pLMN_Identity.size))
							servedPlmnPerCell.PlmnID.Size = int(servedPlmnPerCell_C.pLMN_Identity.size)

							if servedPlmnPerCell_C.du_PM_5GC != nil {
								duPM5GC := &DUPM5GCContainerType{}
								duPM5GC_C := (*C.FGC_DU_PM_Container_t)(servedPlmnPerCell_C.du_PM_5GC)

								duPM5GC.SlicePerPlmnPerCellCount = int(duPM5GC_C.slicePerPlmnPerCellList.list.count)
								for l := 0; l < duPM5GC.SlicePerPlmnPerCellCount; l++ {
									slicePerPlmnPerCell := &duPM5GC.SlicePerPlmnPerCells[l]
									var sizeof_SlicePerPlmnPerCellListItem_t *C.SlicePerPlmnPerCellListItem_t
									slicePerPlmnPerCell_C := *(**C.SlicePerPlmnPerCellListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(duPM5GC_C.slicePerPlmnPerCellList.list.array)) + (uintptr)(l)*unsafe.Sizeof(sizeof_SlicePerPlmnPerCellListItem_t)))

									slicePerPlmnPerCell.SliceID.SST.Buf = C.GoBytes(unsafe.Pointer(slicePerPlmnPerCell_C.sliceID.sST.buf), C.int(slicePerPlmnPerCell_C.sliceID.sST.size))
									slicePerPlmnPerCell.SliceID.SST.Size = int(slicePerPlmnPerCell_C.sliceID.sST.size)

									if slicePerPlmnPerCell_C.sliceID.sD != nil {
										slicePerPlmnPerCell.SliceID.SD = &OctetString{}
										slicePerPlmnPerCell.SliceID.SD.Buf = C.GoBytes(unsafe.Pointer(slicePerPlmnPerCell_C.sliceID.sD.buf), C.int(slicePerPlmnPerCell_C.sliceID.sD.size))
										slicePerPlmnPerCell.SliceID.SD.Size = int(slicePerPlmnPerCell_C.sliceID.sD.size)
									}

									slicePerPlmnPerCell.FQIPERSlicesPerPlmnPerCellCount = int(slicePerPlmnPerCell_C.fQIPERSlicesPerPlmnPerCellList.list.count)
									for m := 0; m < slicePerPlmnPerCell.FQIPERSlicesPerPlmnPerCellCount; m++ {
										fQIPerSlicesPerPlmnPerCell := &slicePerPlmnPerCell.FQIPERSlicesPerPlmnPerCells[m]
										var sizeof_FQIPERSlicesPerPlmnPerCellListItem_t *C.FQIPERSlicesPerPlmnPerCellListItem_t
										fQIPerSlicesPerPlmnPerCell_C := *(**C.FQIPERSlicesPerPlmnPerCellListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(slicePerPlmnPerCell_C.fQIPERSlicesPerPlmnPerCellList.list.array)) + (uintptr)(m)*unsafe.Sizeof(sizeof_FQIPERSlicesPerPlmnPerCellListItem_t)))

										fQIPerSlicesPerPlmnPerCell.FiveQI = int64(fQIPerSlicesPerPlmnPerCell_C.fiveQI)

										if fQIPerSlicesPerPlmnPerCell_C.dl_PRBUsage != nil {
											fQIPerSlicesPerPlmnPerCell.PrbUsage.DL = int64(*fQIPerSlicesPerPlmnPerCell_C.dl_PRBUsage)
										} else {
											fQIPerSlicesPerPlmnPerCell.PrbUsage.DL = -1
										}

										if fQIPerSlicesPerPlmnPerCell_C.ul_PRBUsage != nil {
											fQIPerSlicesPerPlmnPerCell.PrbUsage.UL = int64(*fQIPerSlicesPerPlmnPerCell_C.ul_PRBUsage)
										} else {
											fQIPerSlicesPerPlmnPerCell.PrbUsage.UL = -1
										}
									}
								}

								servedPlmnPerCell.DUPM5GC = duPM5GC
							}

							if servedPlmnPerCell_C.du_PM_EPC != nil {
								duPMEPC := &DUPMEPCContainerType{}
								duPMEPC_C := (*C.EPC_DU_PM_Container_t)(servedPlmnPerCell_C.du_PM_EPC)

								duPMEPC.PerQCIReportCount = int(duPMEPC_C.perQCIReportList.list.count)
								for l := 0; l < duPMEPC.PerQCIReportCount; l++ {
									perQCIReport := &duPMEPC.PerQCIReports[l]
									var sizeof_PerQCIReportListItem_t *C.PerQCIReportListItem_t
									perQCIReport_C := *(**C.PerQCIReportListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(duPMEPC_C.perQCIReportList.list.array)) + (uintptr)(l)*unsafe.Sizeof(sizeof_PerQCIReportListItem_t)))

									perQCIReport.QCI = int64(perQCIReport_C.qci)

									if perQCIReport_C.dl_PRBUsage != nil {
										perQCIReport.PrbUsage.DL = int64(*perQCIReport_C.dl_PRBUsage)
									} else {
										perQCIReport.PrbUsage.DL = -1
									}

									if perQCIReport_C.ul_PRBUsage != nil {
										perQCIReport.PrbUsage.UL = int64(*perQCIReport_C.ul_PRBUsage)
									} else {
										perQCIReport.PrbUsage.UL = -1
									}
								}

								servedPlmnPerCell.DUPMEPC = duPMEPC
							}
						}
					}

					pfContainer.Container = oDU_PF
					//fmt.Println("//// e2sm in type1 pfContainer.Container= oDU_PF %x", pfContainer.Container)
				} else if pfContainer.ContainerType == 2 {
					fmt.Println("//////e2sm if pfContainer.ContainerType == 2")
					oCU_CP_PF := &OCUCPPFContainerType{}
					oCU_CP_PF_C := *(**C.OCUCP_PF_Container_t)(unsafe.Pointer(&pmContainer_C.performanceContainer.choice[0]))

					if oCU_CP_PF_C.gNB_CU_CP_Name != nil {
						oCU_CP_PF.GNBCUCPName = &PrintableString{}
						oCU_CP_PF.GNBCUCPName.Buf = C.GoBytes(unsafe.Pointer(oCU_CP_PF_C.gNB_CU_CP_Name.buf), C.int(oCU_CP_PF_C.gNB_CU_CP_Name.size))
						oCU_CP_PF.GNBCUCPName.Size = int(oCU_CP_PF_C.gNB_CU_CP_Name.size)
					}

					if oCU_CP_PF_C.cu_CP_Resource_Status.numberOfActive_UEs != nil {
						oCU_CP_PF.CUCPResourceStatus.NumberOfActiveUEs = int64(*oCU_CP_PF_C.cu_CP_Resource_Status.numberOfActive_UEs)
					}

					pfContainer.Container = oCU_CP_PF
					fmt.Println("//// e2sm in type2 pfContainer.Container = oCU_CP_PF %x", pfContainer.Container)
				} else if pfContainer.ContainerType == 3 {
					fmt.Println("///////entered e2sm pfcontainer type ==3")
					oCU_UP_PF := &OCUUPPFContainerType{}
					oCU_UP_PF_C := *(**C.OCUUP_PF_Container_t)(unsafe.Pointer(&pmContainer_C.performanceContainer.choice[0]))

					if oCU_UP_PF_C.gNB_CU_UP_Name != nil {
						oCU_UP_PF.GNBCUUPName = &PrintableString{}
						oCU_UP_PF.GNBCUUPName.Buf = C.GoBytes(unsafe.Pointer(oCU_UP_PF_C.gNB_CU_UP_Name.buf), C.int(oCU_UP_PF_C.gNB_CU_UP_Name.size))
						oCU_UP_PF.GNBCUUPName.Size = int(oCU_UP_PF_C.gNB_CU_UP_Name.size)
					}

					oCU_UP_PF.CUUPPFContainerItemCount = int(oCU_UP_PF_C.pf_ContainerList.list.count)
					for j := 0; j < oCU_UP_PF.CUUPPFContainerItemCount; j++ {
						cuUPPFContainer := &oCU_UP_PF.CUUPPFContainerItems[j]
						var sizeof_PF_ContainerListItem_t *C.PF_ContainerListItem_t
						cuUPPFContainer_C := *(**C.PF_ContainerListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(oCU_UP_PF_C.pf_ContainerList.list.array)) + (uintptr)(j)*unsafe.Sizeof(sizeof_PF_ContainerListItem_t)))

						cuUPPFContainer.InterfaceType = int64(cuUPPFContainer_C.interface_type)

						cuUPPFContainer.OCUUPPMContainer.CUUPPlmnCount = int(cuUPPFContainer_C.o_CU_UP_PM_Container.plmnList.list.count)
						for k := 0; k < cuUPPFContainer.OCUUPPMContainer.CUUPPlmnCount; k++ {
							cuUPPlmn := &cuUPPFContainer.OCUUPPMContainer.CUUPPlmns[k]
							var sizeof_PlmnID_List_t *C.PlmnID_List_t
							cuUPPlmn_C := *(**C.PlmnID_List_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(cuUPPFContainer_C.o_CU_UP_PM_Container.plmnList.list.array)) + (uintptr)(k)*unsafe.Sizeof(sizeof_PlmnID_List_t)))

							cuUPPlmn.PlmnID.Buf = C.GoBytes(unsafe.Pointer(cuUPPlmn_C.pLMN_Identity.buf), C.int(cuUPPlmn_C.pLMN_Identity.size))
							cuUPPlmn.PlmnID.Size = int(cuUPPlmn_C.pLMN_Identity.size)

							if cuUPPlmn_C.cu_UP_PM_5GC != nil {
								cuUPPM5GC := &CUUPPM5GCType{}
								cuUPPM5GC_C := (*C.FGC_CUUP_PM_Format_t)(cuUPPlmn_C.cu_UP_PM_5GC)

								cuUPPM5GC.SliceToReportCount = int(cuUPPM5GC_C.sliceToReportList.list.count)
								for l := 0; l < cuUPPM5GC.SliceToReportCount; l++ {
									sliceToReport := &cuUPPM5GC.SliceToReports[l]
									var sizeof_SliceToReportListItem_t *C.SliceToReportListItem_t
									sliceToReport_C := *(**C.SliceToReportListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(cuUPPM5GC_C.sliceToReportList.list.array)) + (uintptr)(l)*unsafe.Sizeof(sizeof_SliceToReportListItem_t)))

									sliceToReport.SliceID.SST.Buf = C.GoBytes(unsafe.Pointer(sliceToReport_C.sliceID.sST.buf), C.int(sliceToReport_C.sliceID.sST.size))
									sliceToReport.SliceID.SST.Size = int(sliceToReport_C.sliceID.sST.size)

									if sliceToReport_C.sliceID.sD != nil {
										sliceToReport.SliceID.SD = &OctetString{}
										sliceToReport.SliceID.SD.Buf = C.GoBytes(unsafe.Pointer(sliceToReport_C.sliceID.sD.buf), C.int(sliceToReport_C.sliceID.sD.size))
										sliceToReport.SliceID.SD.Size = int(sliceToReport_C.sliceID.sD.size)
									}

									sliceToReport.FQIPERSlicesPerPlmnCount = int(sliceToReport_C.fQIPERSlicesPerPlmnList.list.count)
									for m := 0; m < sliceToReport.FQIPERSlicesPerPlmnCount; m++ {
										fQIPerSlicesPerPlmn := &sliceToReport.FQIPERSlicesPerPlmns[m]
										var sizeof_FQIPERSlicesPerPlmnListItem_t *C.FQIPERSlicesPerPlmnListItem_t
										fQIPerSlicesPerPlmn_C := *(**C.FQIPERSlicesPerPlmnListItem_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(sliceToReport_C.fQIPERSlicesPerPlmnList.list.array)) + (uintptr)(m)*unsafe.Sizeof(sizeof_FQIPERSlicesPerPlmnListItem_t)))

										fQIPerSlicesPerPlmn.FiveQI = int64(fQIPerSlicesPerPlmn_C.fiveQI)

										if fQIPerSlicesPerPlmn_C.pDCPBytesDL != nil {
											fQIPerSlicesPerPlmn.PDCPBytesDL = &Integer{}
											fQIPerSlicesPerPlmn.PDCPBytesDL.Buf = C.GoBytes(unsafe.Pointer(fQIPerSlicesPerPlmn_C.pDCPBytesDL.buf), C.int(fQIPerSlicesPerPlmn_C.pDCPBytesDL.size))
											fQIPerSlicesPerPlmn.PDCPBytesDL.Size = int(fQIPerSlicesPerPlmn_C.pDCPBytesDL.size)
										}

										if fQIPerSlicesPerPlmn_C.pDCPBytesUL != nil {
											fQIPerSlicesPerPlmn.PDCPBytesUL = &Integer{}
											fQIPerSlicesPerPlmn.PDCPBytesUL.Buf = C.GoBytes(unsafe.Pointer(fQIPerSlicesPerPlmn_C.pDCPBytesUL.buf), C.int(fQIPerSlicesPerPlmn_C.pDCPBytesUL.size))
											fQIPerSlicesPerPlmn.PDCPBytesUL.Size = int(fQIPerSlicesPerPlmn_C.pDCPBytesUL.size)
										}
									}
								}

								cuUPPlmn.CUUPPM5GC = cuUPPM5GC
							}

							if cuUPPlmn_C.cu_UP_PM_EPC != nil {
								cuUPPMEPC := &CUUPPMEPCType{}
								cuUPPMEPC_C := (*C.EPC_CUUP_PM_Format_t)(cuUPPlmn_C.cu_UP_PM_EPC)

								cuUPPMEPC.CUUPPMEPCPerQCIReportCount = int(cuUPPMEPC_C.perQCIReportList.list.count)
								for l := 0; l < cuUPPMEPC.CUUPPMEPCPerQCIReportCount; l++ {
									perQCIReport := &cuUPPMEPC.CUUPPMEPCPerQCIReports[l]
									var sizeof_PerQCIReportListItemFormat_t *C.PerQCIReportListItemFormat_t
									perQCIReport_C := *(**C.PerQCIReportListItemFormat_t)(unsafe.Pointer((uintptr)(unsafe.Pointer(cuUPPMEPC_C.perQCIReportList.list.array)) + (uintptr)(l)*unsafe.Sizeof(sizeof_PerQCIReportListItemFormat_t)))

									perQCIReport.QCI = int64(perQCIReport_C.qci)

									if perQCIReport_C.pDCPBytesDL != nil {
										perQCIReport.PDCPBytesDL = &Integer{}
										perQCIReport.PDCPBytesDL.Buf = C.GoBytes(unsafe.Pointer(perQCIReport_C.pDCPBytesDL.buf), C.int(perQCIReport_C.pDCPBytesDL.size))
										perQCIReport.PDCPBytesDL.Size = int(perQCIReport_C.pDCPBytesDL.size)
									}

									if perQCIReport_C.pDCPBytesUL != nil {
										perQCIReport.PDCPBytesUL = &Integer{}
										perQCIReport.PDCPBytesUL.Buf = C.GoBytes(unsafe.Pointer(perQCIReport_C.pDCPBytesUL.buf), C.int(perQCIReport_C.pDCPBytesUL.size))
										perQCIReport.PDCPBytesUL.Size = int(perQCIReport_C.pDCPBytesUL.size)
									}
								}

								cuUPPlmn.CUUPPMEPC = cuUPPMEPC
							}
						}
					}

					pfContainer.Container = oCU_UP_PF
					fmt.Println("//// e2sm in type3 pfContainer.Container= oCU_UP_PF %x", pfContainer.Container)
				} else {
					fmt.Println("//////e2sm in else Unknown PF Container type indMsg %x", indMsg)
					return indMsg, errors.New("Unknown PF Container type")
				}

				pmContainer.PFContainer = pfContainer
				fmt.Println("/////e2sm after else pmContainer.PFContainer = pfContainer %x", pmContainer.PFContainer)
			}

		}

		indMsg.IndMsg = indMsgFormat1
		fmt.Println("/////e2sm before second else indMsg.IndMsg = indMsgFormat1 %x", indMsg.IndMsg)
	} else {
		fmt.Println("//////e2sm in second else Unknown PF Container type indMsg %x", indMsg)
		return indMsg, errors.New("Unknown RIC Indication Message Format")
	}

	return
}

func (c *E2sm) ParseNRCGI(nRCGI NRCGIType) (CellID string, err error) {
	var plmnID OctetString
	var nrCellID BitString

	//plmnID = nRCGI.PlmnID
	//CellID, _ = c.ParsePLMNIdentity(plmnID.Buf, plmnID.Size)
	
	plmnID = nRCGI.PlmnID
	fmt.Println("plmnID = nRCGI.PlmnID in e2sm parsenrcgi func: %d", plmnID)
	CellID, _ = c.ParsePLMNIdentity(plmnID.Buf, plmnID.Size)
	fmt.Println("CellID in e2sm parsenrcgi func: %d", CellID)

	nrCellID = nRCGI.NRCellID
	fmt.Println("nrCellID = nRCGI.NRCellID in e2sm parsenrcgi func:", nrCellID)
	
	fmt.Println("plmnID.Size= %d", plmnID.Size)
	fmt.Println("nrCellID.Size= %d", nrCellID.Size)

	if plmnID.Size != 3 || nrCellID.Size != 5 {
		return "", errors.New("Invalid input: illegal length of NRCGI")
	}

	var former []uint8 = make([]uint8, 3)
	var latter []uint8 = make([]uint8, 6)

	former[0] = nrCellID.Buf[0] >> 4
	former[1] = nrCellID.Buf[0] & 0xf
	former[2] = nrCellID.Buf[1] >> 4
	latter[0] = nrCellID.Buf[1] & 0xf
	latter[1] = nrCellID.Buf[2] >> 4
	latter[2] = nrCellID.Buf[2] & 0xf
	latter[3] = nrCellID.Buf[3] >> 4
	latter[4] = nrCellID.Buf[3] & 0xf
	latter[5] = nrCellID.Buf[4] >> uint(nrCellID.BitsUnused)

	CellID = CellID + strconv.Itoa(int(former[0])) + strconv.Itoa(int(former[1])) + strconv.Itoa(int(former[2])) + strconv.Itoa(int(latter[0])) + strconv.Itoa(int(latter[1])) + strconv.Itoa(int(latter[2])) + strconv.Itoa(int(latter[3])) + strconv.Itoa(int(latter[4])) + strconv.Itoa(int(latter[5]))
	
	fmt.Println("CellID at the end of parsenrcgi func in e2sm: %d", CellID)

	return
}

func (c *E2sm) ParsePLMNIdentity(buffer []byte, size int) (PlmnID string, err error) {
	if size != 3 {
		fmt.Println("////e2sm entered ParsePLMNIdentity if size != 3")
		return "", errors.New("Invalid input: illegal length of PlmnID")
	}

	var mcc []uint8 = make([]uint8, 3)
	var mnc []uint8 = make([]uint8, 3)

	mcc[0] = buffer[0] >> 4
	mcc[1] = buffer[0] & 0xf
	mcc[2] = buffer[1] >> 4
	mnc[0] = buffer[1] & 0xf
	mnc[1] = buffer[2] >> 4
	mnc[2] = buffer[2] & 0xf

	if mnc[0] == 0xf {
		PlmnID = strconv.Itoa(int(mcc[0])) + strconv.Itoa(int(mcc[1])) + strconv.Itoa(int(mcc[2])) + strconv.Itoa(int(mnc[1])) + strconv.Itoa(int(mnc[2]))
	} else {
		PlmnID = strconv.Itoa(int(mcc[0])) + strconv.Itoa(int(mcc[1])) + strconv.Itoa(int(mcc[2])) + strconv.Itoa(int(mnc[0])) + strconv.Itoa(int(mnc[1])) + strconv.Itoa(int(mnc[2]))
	}

	return
}

func (c *E2sm) ParseSliceID(sliceID SliceIDType) (combined int32, err error) {
	if sliceID.SST.Size != 1 || (sliceID.SD != nil && sliceID.SD.Size != 3) {
		return 0, errors.New("Invalid input: illegal length of sliceID")
	}

	var temp uint8
	var sst int32
	var sd int32

	byteBuffer := bytes.NewBuffer(sliceID.SST.Buf)
	binary.Read(byteBuffer, binary.BigEndian, &temp)
	sst = int32(temp)

	if sliceID.SD == nil {
		combined = sst << 24
	} else {
		for i := 0; i < sliceID.SD.Size; i++ {
			byteBuffer = bytes.NewBuffer(sliceID.SD.Buf[i : i+1])
			binary.Read(byteBuffer, binary.BigEndian, &temp)
			sd = sd*256 + int32(temp)
		}
		combined = sst<<24 + sd
	}

	return
}

func (c *E2sm) ParseInteger(buffer []byte, size int) (value int64, err error) {
	var temp uint8
	var byteBuffer *bytes.Buffer

	for i := 0; i < size; i++ {
		byteBuffer = bytes.NewBuffer(buffer[i : i+1])
		binary.Read(byteBuffer, binary.BigEndian, &temp)
		value = value*256 + int64(temp)
	}

	return
}

func (c *E2sm) ParseTimestamp(buffer []byte, size int) (timestamp *Timestamp, err error) {
	var temp uint8
	var byteBuffer *bytes.Buffer
	var index int
	var sec int64
	var nsec int64

	for index := 0; index < size-8; index++ {
		byteBuffer = bytes.NewBuffer(buffer[index : index+1])
		binary.Read(byteBuffer, binary.BigEndian, &temp)
		sec = sec*256 + int64(temp)
	}

	for index = size - 8; index < size; index++ {
		byteBuffer = bytes.NewBuffer(buffer[index : index+1])
		binary.Read(byteBuffer, binary.BigEndian, &temp)
		nsec = nsec*256 + int64(temp)
	}

	timestamp = &Timestamp{TVsec: sec, TVnsec: nsec}
	return
}
