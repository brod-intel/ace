//Copyright 2019, Intel Corporation

package main

import (
	"fmt"
	"helpers"
	"math/rand"
	"memberupdatex"
	"os"
	"strconv"
	"strings"
	"time"
)

//Global variable to be used across the file
var (
	acelog *helpers.Logger
)

// InitHandler intialises the handler
func InitHandler() error {
	acelog = helpers.GetLogger()
	return nil
}

// leaderState will check if check if leader is in progress and
// accordingly set or delete the tag "waitingforleader".
func leaderState() error {

	//Read the hostname , add the error code
	hostname, _ := os.Hostname()
	acelog.Debug("hostname: ", hostname)
	acelog.Debug("len of hostname: ", len(hostname))

	acelog.Info("Running Serf-member Update")

	tag := make(map[string]string)
	tag["role"] = "leader"
	tag["inprocess"] = "true"

	leaderInProcess, err := helpers.ListMemberByTags(tag)
	if err != nil {
		acelog.Error("Value of err in getting Leader name: ", err)
	}

	acelog.Debug("leader inprocess: ", leaderInProcess)

	if leaderInProcess != "" {
		acelog.Debug("value of leaderinprogress", leaderInProcess)
		if leaderInProcess != hostname {

			serfTagsOp := helpers.SetWaitingForLeaderTag("true")
			if serfTagsOp != nil {
				acelog.Error("Error SetWaitingForLeaderTag :", serfTagsOp)
			}
			return fmt.Errorf("SetWaitingForLeaderTag")
		}
		return nil
	}

	tags := make(map[string]string)
	tags["waitingforleader"] = "true"

	checkLeaderStatus, err := helpers.MemberNameByTagsAndName(tags, hostname)
	acelog.Debug("value of getLeaderStatus ", checkLeaderStatus)

	if err != nil {
		return fmt.Errorf("Error to checkleader Status %v", err)

	}

	if strings.Contains(checkLeaderStatus, hostname) {
		acelog.Debug("deleting the tag \"waitingforleader\"")

		serfTagsOp := helpers.DeleteSerfTag("waitingforleader")
		if serfTagsOp != nil {
			acelog.Error("Error Deleting, WaitingForLeaderTag :", serfTagsOp)
		}
		return fmt.Errorf("deleting the tag waitingforleader")
	}
	return nil
}

// WorkerState will check if check if worker is in progress and
// accordingly set or delete the tag "waitingfoacerker".
func workerState() error {

	hostname, _ := os.Hostname()
	acelog.Debug("hostname:  ", hostname)
	acelog.Debug("len of hostname: ", len(hostname))

	acelog.Debug("Running Serf-member Update")

	tag := make(map[string]string)
	tag["role"] = "worker"
	tag["inprocess"] = "true"

	workerInProcess, err := helpers.ListMemberByTags(tag)
	if err != nil {
		return fmt.Errorf("%v", err)

	}

	acelog.Debug("worker inprocess: ", workerInProcess)

	if workerInProcess != "" {
		acelog.Debug("value of workerinprogress", workerInProcess)
		if workerInProcess != hostname {

			serfTagsOp := helpers.SetWaitingForWorkerTag("true")
			if serfTagsOp != nil {
				acelog.Error("Error SetWaitingForWorkerTag :", serfTagsOp)
			}
			return fmt.Errorf("SetWaitingForWorkerTag")
		}
		return nil
	}

	tags := make(map[string]string)
	tags["waitingfoacerker"] = "true"

	checkworkerStatus, err := helpers.MemberNameByTagsAndName(tags, hostname)
	acelog.Debug("value of getworkerStatus ", checkworkerStatus)

	if err != nil {
		return fmt.Errorf("Error code %v", err)
	}

	if strings.Contains(checkworkerStatus, hostname) {
		acelog.Debug("deleting the tag \"waitingfoacerker\"")

		serfTagsOp := helpers.DeleteSerfTag("waitingfoacerker")
		if serfTagsOp != nil {
			acelog.Error("Error Deleting, WaitingFoacerkerTag :", serfTagsOp)
		}
		return fmt.Errorf("deleting the tag waitingfoacerker")
	}
	return nil
}

// memberUpdateCheckAndSetTag Set the tag to worker/leader
//If Count =1, set leader
//else, If leader exists; set worker else leader
func memberUpdateCheckAndSetTag() error {

	countAlive := helpers.CountAliveMembers()
	acelog.Debug("CountAliveMembers :", countAlive)

	if countAlive == 1 {
		acelog.Debug("Empty role tag, setting 'leader' tag")

		serfTagsOp := helpers.SetRoleTag("leader")
		if serfTagsOp != nil {
			acelog.Error("Error in updating tag, Role ", serfTagsOp)
			return serfTagsOp
		}

	} else {

		tags := make(map[string]string)
		tags["role"] = "leader"

		aliveLeaderStatus, err := helpers.MemberIPByTagsAndStatus(tags, "alive")
		if err != nil {
			acelog.Error("helpers.MemberIPByTagsAndStatus Failed: ", err)
			return err
		}

		acelog.Debug("alive LeaderStatus: ", aliveLeaderStatus)

		if len(aliveLeaderStatus) != 0 {
			acelog.Debug("Empty role tag, leader exists setting 'worker' tag")

			//Setting the role tag "worker", as Leader exists
			serfTagsOp := helpers.SetRoleTag("worker")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag, Role ", serfTagsOp)
				return serfTagsOp
			}

		} else {

			acelog.Debug("Empty role tag, leader does not exist, setting 'leader' tag")

			//No leader exists, so assigning leader tag.
			serfTagsOp := helpers.SetRoleTag("leader")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag, Role ", serfTagsOp)
				return serfTagsOp
			}

		}

	}
	return nil
}

//maintainQuorum: leader/reachable node not allowed
//Worker Node with the minimum value of init will be promoted as manager.
func maintainQuorum() error {

	inspectManagerStatus, err := helpers.GetNodeStatus("reachable")
	if err != nil {
		acelog.Error("Failed to check the node status", err)
		return err
	}

	acelog.Debug("status of ManagerStatus.Reachable:", inspectManagerStatus)

	if inspectManagerStatus == true {
		acelog.Debug("reachable role returning from maintainQuorum")
		return nil
	}

	inspectManagerStatus, err = helpers.GetNodeStatus("leader")
	if err != nil {
		acelog.Error("Failed to check the node status")
		return err
	}

	acelog.Debug("status of ManagerStatus.leader:", inspectManagerStatus)

	if inspectManagerStatus == true {
		acelog.Debug("leader role returning from maintainQuorum")
		return nil
	}

	countAlive := helpers.CountAliveMembers()
	acelog.Debug("CountAliveMembers :", countAlive)

	reachableCount, err := helpers.SerfQuery("godocker", "reachable")
	if err != nil {
		acelog.Error("Node Failed to check the ReachableNodeCount :\n")
		return err
	}

	countReachableNodes, _ := strconv.Atoi(reachableCount)

	acelog.Debug("CountReachableNodes:", reachableCount)

	//Condition to check if ACE needs promotion of any node.
	if countAlive > 3 && countReachableNodes < 2 {

		//Check whose init is first, that node is eligible to be promoted.
		swarmID, err := helpers.GetNodeIDForWorkerWithMinTagValue()
		if err != nil {
			acelog.Error("Failed to get worker with Min Tag Value ", err)
			return err
		}

		if swarmID == "" {
			acelog.Debug("Not eligible for promotion")
			return nil
		}

		//Node could be a part of previous swarm, Leave and rejoin
		time.Sleep(5 * time.Second)
		err = helpers.SwarmLeave(true)
		time.Sleep(10 * time.Second)
		if err != nil {
			return fmt.Errorf("Failed to leave the node %v", err)
		}

		tags := make(map[string]string)
		tags["role"] = "leader"

		serfLeader, err := helpers.MemberIPByTagsAndStatus(tags, "alive")
		if err != nil {
			return fmt.Errorf("MemberIPByTagsAndStatus Failed, Value of serfLeader %v", string(serfLeader))
		}

		serfAdvertiseIface, err := helpers.GetIPAddr()
		if err != nil {
			acelog.Error("Failed to get network IP\n")
			return err
		}

		if len(serfAdvertiseIface) == 0 {
			return fmt.Errorf("Error IP is not set ")
		}

		//To handle Network Flapping
		retry := 1
		for retry < 4 {

			err = memberupdatex.JoinSwarm(serfAdvertiseIface, serfLeader, "manager")
			time.Sleep(10 * time.Second)

			if err == nil {
				break
			}

			if err != nil && retry >= 4 {
				serfTagsOp := helpers.SetTempTag()
				if serfTagsOp != nil {
					acelog.Error("Error in update tag, Inprocess ", serfTagsOp)
					return serfTagsOp
				}

				serfTagsOp = helpers.DeleteSerfTag("temp")
				if serfTagsOp != nil {
					acelog.Error("Error of in Deleting Tag, gluster: ", serfTagsOp)
					return serfTagsOp
				}

				return fmt.Errorf("Node failed to join %v", swarmID)
			}

			retry++
			time.Sleep(5 * time.Second) // To wait for leader node to update token
		}

		var removeDownNodes []string
		retry = 1
		for retry < 10 {

			removeDownNodes, err = helpers.GetNodeIDByState("down")
			if err != nil {
				return fmt.Errorf("Node Status for down failed ")
			}

			if len(removeDownNodes) != 0 {
				break
			}
			retry++
			time.Sleep(5 * time.Second) // To wait for down node
		}

		for _, removeNode := range removeDownNodes {

			if len(removeNode) != 0 {

				err := helpers.DemoteNode(removeNode)
				if err != nil {
					acelog.Error("Failed to demote the node ", err.Error())
				}

				err = helpers.RemoveNode(removeNode, false)
				if err != nil {
					acelog.Error("Failed to remove the node", err.Error())
					return err
				}
			}
		}

		acelog.Debug("Node Joined successfully", swarmID)

		err = memberupdatex.SetSwarmTag()
		if err != nil {
			acelog.Error("Failed to set the tag, setSwarmTag", err.Error())
			return err
		}
	}

	return nil
}

func manageSwarmStatus(swarmManager bool, statusSwarm string) error {

	inspectManagerStatus, err := helpers.GetNodeStatus("reachable")
	if err != nil {
		acelog.Debug("Failed to check the node status")
		return err
	}

	acelog.Debug("status of ManagerStatus.Reachable:", inspectManagerStatus)

	if inspectManagerStatus == true {
		acelog.Debug("Worker role, there are less than or equal to 3 serf members and connected to Docker Swarm as a manager")

	} else {

		if len(statusSwarm) != 0 && swarmManager == false {

			serfTagsOp := helpers.SetSwarmTag("")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag, Swarm ", serfTagsOp)
				return serfTagsOp
			}

			return fmt.Errorf("Worker role, there are less than or equal to 3 serf members, swarm tag populated and not a swarm manager.  Resetting swarm tag")
		}

	}

	return nil
}

// memberUpdateWorker will update the tag as per the status
// If the only node, update the tag to leader
// If the only node, and not a swarm manager , set leader tag & clear swarm tag.

func memberUpdateWorker() error {

	acelog.Debug("In MemberUpdateWorker\n")

	inspectManagerStatus, err := helpers.GetNodeStatus("leader")
	if err != nil {
		acelog.Debug("Failed to check the node status")
		return err
	}

	swarmManager, err := helpers.CheckIfManager()
	if err != nil {
		acelog.Error("Failed to check manager status")
		return err
	}

	countAlive := helpers.CountAliveMembers()
	acelog.Debug("CountAliveMembers :", countAlive)

	countTotal := helpers.CountTotalMembers()
	acelog.Debug("CountTotalMembers :", countTotal)

	if countAlive < countTotal {
		return fmt.Errorf("One or more members has left, skipping member-update to allow member-failed perform the action")
	}

	statusSwarm, _ := os.LookupEnv("SERF_TAG_SWARM")
	acelog.Debug("value of statusSwarm ", statusSwarm)

	if countAlive == 1 && len(statusSwarm) != 0 && inspectManagerStatus == true {

		serfTagsOp := helpers.SetRoleTag("leader")
		if serfTagsOp != nil {
			acelog.Error("Error in updating tag, Role ", serfTagsOp)
			return serfTagsOp
		}

		return fmt.Errorf("worker role, only 1 serf member, swarm tag exists, a Swarm Manager 'leader', setting 'leader' tag")
	}

	if countAlive == 1 && len(statusSwarm) != 0 && swarmManager == false {

		serfTagsOp := helpers.SetSwarmTag("")
		if serfTagsOp != nil {
			acelog.Error("Error in updating tag, Swarm ", serfTagsOp)
		}

		serfTagsOp = helpers.SetRoleTag("leader")
		if serfTagsOp != nil {
			acelog.Error("Error in updating tag, Role ", serfTagsOp)
		}
		return fmt.Errorf("worker role, only 1 serf member, swarm tag exists, but not a Swarm Manager, 'leader' tag set and 'swarm' tag cleared")
	}

	tags := make(map[string]string)
	tags["role"] = "leader"

	aliveLeaderStatus, err := helpers.MemberIPByTagsAndStatus(tags, "alive")
	if err != nil {
		acelog.Debug("MemberIPByTagsAndStatus Failed:  ", aliveLeaderStatus)
		return err
	}

	acelog.Debug("alive Leader status: ", aliveLeaderStatus)

	if len(aliveLeaderStatus) != 0 {
		acelog.Debug("A leader exists, no need to reset")
	} else {

		if inspectManagerStatus == true {

			serfTagsOp := helpers.SetRoleTag("leader")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag ", serfTagsOp)
			}

			return fmt.Errorf("Worker role and Swarm Manager 'leader', update role tag to 'leader' ")
		}
	}

	if countAlive <= 3 {

		//if worker node is not reachable & part of swarm, clear the swarm tag
		err := manageSwarmStatus(swarmManager, statusSwarm)
		if err != nil {
			return err
		}
		return nil
	}

	statusGluster, _ := os.LookupEnv("SERF_TAG_GLUSTER")
	acelog.Debug("value of statusGluster:", statusGluster)

	//Only if gluster is set, we are eligible to go to swarm
	if len(statusGluster) != 0 {

		acelog.Debug("Maintaining Quorum")
		err = maintainQuorum()
		if err != nil {
			acelog.Error("maintainQuorum Failed ", err)
			return err
		}
	}

	return nil

}

// func memberUpdateLeader will update the tag as per the status
// If ManagerStatus.Leader is false, set the node as worker
// If not a  swarm manager, clear the swarm tag.

func memberUpdateLeader() error {

	countAlive := helpers.CountAliveMembers()
	acelog.Debug("CountAliveMembers :", countAlive)

	countTotal := helpers.CountTotalMembers()
	acelog.Debug("CountTotalMembers :", countTotal)

	if countAlive < countTotal {
		return fmt.Errorf("one or more members has left, skipping member-update to allow member-failed perform the action")
	}

	statusSwarm, _ := os.LookupEnv("SERF_TAG_SWARM")

	acelog.Debug("value of statusSwarm:", string(statusSwarm))

	if len(statusSwarm) != 0 {

		swarmManager, err := helpers.CheckIfManager()
		if err != nil {
			acelog.Error("Failed to check manager status")
			return err
		}

		if countAlive == 1 && swarmManager == false {
			acelog.Debug("Leader role but not a Swarm Manager, clearing swarm tag")

			serfTagsOp := helpers.SetSwarmTag("")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag ", serfTagsOp)
			}

			return fmt.Errorf("leader role but not a swarm manager. cleared swarm tag")
		}

		inspectManagerStatus, err := helpers.GetNodeStatus("leader")
		if err != nil {
			acelog.Error("Failed to check the node status")
			return err
		}

		if inspectManagerStatus == true {
			acelog.Debug("Leader role and Swarm Manager 'leader'")
		} else {
			acelog.Debug("Leader role but not a Swarm Manager 'leader', setting 'worker' tag")

			serfTagsOp := helpers.SetRoleTag("worker")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag ", serfTagsOp)
			}
		}

	} else {

		//Sort by init and return the hostname
		tags := make(map[string]string)
		tags["role"] = "leader"

		leaderName, err := helpers.MemberNameByTagsAndStatus(tags, "alive")
		if err != nil {
			acelog.Debug("MemberNameByTagsAndStatus Failed, Value of serfLeader: ", string(leaderName))
			return err
		}

		acelog.Debug("value of leadername: ", leaderName)

		serfSelfName, _ := os.LookupEnv("SERF_SELF_NAME")
		acelog.Debug("value of SerfSelfName : ", serfSelfName)

		if leaderName != serfSelfName {

			serfTagsOp := helpers.SetRoleTag("worker")
			if serfTagsOp != nil {
				acelog.Error("Error in updating tag ", serfTagsOp)
			}
			return fmt.Errorf("Leader role and swarm tag is empty, Serf Leader does not match this system's name so setting 'worker' tag")
		}

	}

	return nil
}

func main() {

	// Initialise handler
	InitHandler()

	//check for network status
	if helpers.CheckNetworkStatus() != nil {
		return
	}

	rand.Seed(time.Now().UnixNano()) //Seed for random sleep

	//Check if leaderInprocess, then set/delete the tag and exit
	status := leaderState()
	if status != nil {
		acelog.Error(status.Error())
		return
	}

	//Check if workerInprocess, then set/delete the tag and exit
	status = workerState()
	if status != nil {
		acelog.Error(status.Error())
		return
	}

	//Debug Prints
	statusRole, _ := os.LookupEnv("SERF_TAG_ROLE")
	acelog.Debug("value of statusRole", statusRole)

	if !helpers.IsValidRole(statusRole) {
		acelog.Debug("SERF Tag from Env is invalid. ", statusRole)
		statusRole = ""
	}

	if len(statusRole) == 0 {
		err := memberUpdateCheckAndSetTag()
		if err != nil {
			acelog.Error("memberUpdateCheckAndSetTag Failed", err.Error())
			return
		}

	} else if statusRole == "leader" {

		status = memberUpdateLeader()
		if status != nil {
			acelog.Error(status.Error())
			return
		}

		acelog.Debug("memberUpdateLeader Successfully executed")

	} else if statusRole == "worker" {

		status = memberUpdateWorker()
		if status != nil {
			acelog.Error(status.Error())
			return
		}

		acelog.Debug("memberUpdateWorker Successfully executed")

	}

	statusGluster, _ := os.LookupEnv("SERF_TAG_GLUSTER")
	acelog.Debug("value of statusGluster:", statusGluster)
	acelog.Debug("len of statusGluster:", len(statusGluster))

	statusSwarm, _ := os.LookupEnv("SERF_TAG_SWARM")
	acelog.Debug("value of statusSwarm ", statusSwarm)
	acelog.Debug("len of statusSwarm ", len(statusSwarm))

	//Full path is given temporary, a function will be called when go is developed
	if len(statusGluster) == 0 {
		for {
			acelog.Debug("Executing member-update.x  Gluster")
			err := memberupdatex.Gluster()
			if err != nil {
				acelog.Debug("Error in Executing Gluster ", err.Error())
				acelog.Debug("Retrying .......")
			} else {
				acelog.Debug("Completed member-update.x  Gluster")
				break
			}
		}

	} else if len(statusSwarm) == 0 { //Execute the serf swarm who will take care of promotion/demotion
		for {
			acelog.Debug("Executing member-update.x  swarm")
			err := memberupdatex.Swarm()
			if err != nil {
				acelog.Error("Error in Executing swarm ", err.Error())
				acelog.Error("Retrying .......")
			} else {
				acelog.Debug("Completed member-update.x  swarm")
				break
			}
		}

	} else {
		for n := 1; n < 30; n++ {
			acelog.Debug("Executing member-update.x arbiter")
			err := memberupdatex.Arbiter()
			if err != nil {
				acelog.Error("Error in Executing arbiter ", err.Error())
				acelog.Error("Retrying .......")
				time.Sleep(1 * time.Second)
			} else {
				acelog.Debug("Completed member-update.x  arbiter")
				break
			}
		}
	}
}
