package stat

import (
	"reflect"
	"runtime"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/model"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"
	
)

//Controller
type SrvStatistics_Controller struct {
	gobizap.Base_Controller
}

func NewSrvStatistics_Controller() *SrvStatistics_Controller{
	c := &SrvStatistics_Controller{gobizap.Base_Controller{ID: "SrvStatistics", PublicMethods: make(gobizap.PublicMethodCollection)}}
	
	//************************** method get_statistics **********************************
	c.PublicMethods["get_statistics"] = &SrvStatistics_Controller_get_statistics{
		gobizap.Base_PublicMethod{
			ID: "get_statistics",
		},
	}
	
	return c
}

//**************************************************************************************
//Public method: get_statistics
type SrvStatistics_Controller_get_statistics struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *SrvStatistics_Controller_get_statistics) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	return res, nil
}

type statServer interface{
	GetStatistics() SrvStater
}

//custom method
func (pm *SrvStatistics_Controller_get_statistics) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	/*rows := make([]model.ModelRow, 0)
	for srv_name, srv := range app.GetServers() {
		if stat_srv, ok := srv.(statServer); ok {
			stat := stat_srv.GetStatistics()
			m_row := &SrvStatistics{Name: fields.ValText{TypedValue: srv_name},
				Max_client_count: fields.ValInt{TypedValue: int64(stat.GetMaxClientCount())},
				Client_count: fields.ValInt{TypedValue: int64(stat.GetClientCount())},
				Downloaded_bytes: fields.ValUint{TypedValue: stat.GetDownloadedBytes()},
				Uploaded_bytes: fields.ValUint{TypedValue: stat.GetUploadedBytes()},
				Handshakes: fields.ValUint{TypedValue: stat.GetHandshakes()},
				Run_seconds: fields.ValUint{TypedValue: stat.GetRunSeconds()},				
			}
			rows = append(rows, m_row)
		}
	}
	resp.AddModel(&model.Model{ID: model.ModelID("SrvStatistics_Model"), Rows: rows})
	*/
	//mem stat
	var mem_stat runtime.MemStats
        runtime.ReadMemStats(&mem_stat)
        
        ms_rows := make([]model.ModelRow, 1)	
        ms_rows[0] = &MemStat{Sys: fields.ValUint{TypedValue: mem_stat.Sys},
		Lookups: fields.ValUint{TypedValue: mem_stat.Lookups},		
		Heap_sys: fields.ValUint{TypedValue: mem_stat.HeapSys},
		Heap_inuse: fields.ValUint{TypedValue: mem_stat.HeapInuse},
		Heap_objects: fields.ValUint{TypedValue: mem_stat.HeapObjects},
		Stack_sys: fields.ValUint{TypedValue: mem_stat.StackSys},
		Stack_inuse: fields.ValUint{TypedValue: mem_stat.StackInuse},        
        }
	resp.AddModel(&model.Model{ID: model.ModelID("MemStat_Model"), Rows: ms_rows})
	
	return nil
}

