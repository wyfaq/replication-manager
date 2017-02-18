package cluster

import (
	"time"

	"github.com/tanji/replication-manager/dbhelper"
)

func (cluster *Cluster) testSwitchOverLongTransactionWithoutCommitNoRplCheckNoSemiSync() bool {
	cluster.conf.RplChecks = false
	cluster.conf.MaxDelay = 8
	cluster.LogPrintf("TESTING : Starting Test %s", "testSwitchOverLongTransactionNoRplCheckNoSemiSync")
	for _, s := range cluster.servers {
		_, err := s.Conn.Exec("set global rpl_semi_sync_master_enabled='OFF'")
		if err != nil {
			cluster.LogPrintf("TESTING : %s", err)
		}
		_, err = s.Conn.Exec("set global rpl_semi_sync_slave_enabled='OFF'")
		if err != nil {
			cluster.LogPrintf("TESTING : %s", err)
		}
	}

	SaveMasterURL := cluster.master.URL
	go dbhelper.InjectLongTrx(cluster.master.Conn, 20)
	for i := 0; i < 1; i++ {

		cluster.LogPrintf("INFO :  Master is %s", cluster.master.URL)

		switchoverChan <- true

		cluster.waitFailoverEnd()
		cluster.LogPrintf("INFO : New Master  %s ", cluster.master.URL)

	}
	for _, s := range cluster.slaves {
		dbhelper.StartSlave(s.Conn)
	}
	time.Sleep(2 * time.Second)
	if cluster.master.URL != SaveMasterURL {
		cluster.LogPrintf("INFO : Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, cluster.master.URL)
		return false
	}
	return true
}
