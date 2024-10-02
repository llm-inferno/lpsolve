/*********************************************
 * OPL 22.1.1.0 Model
 * Author: tantawi
 * Creation Date: Sep 3, 2024 at 4:01:58 PM
 *********************************************/

using CPLEX;
 
int numServers = ...;
int	numAccelerators = ...;
int numVars = numServers * numAccelerators;

range servers = 0..numServers-1;
range accelerators = 0..numAccelerators-1;
range vars = 0..numVars-1;

float instanceCost[accelerators] = ...;

int numInstancesPerReplica[servers][accelerators] = ...;
float ratePerReplica[servers][accelerators] = ...;
float arrivalRates[servers] = ...;

float costVector[vars];
float rateVector[servers][vars];
int excluded[vars];
execute {
  for(var i in servers) {
    for(var j in accelerators) {
      costVector[i * numAccelerators + j] = numInstancesPerReplica[i][j] * instanceCost[j]
      rateVector[i][i * numAccelerators + j] = ratePerReplica[i][j]
      if (ratePerReplica[i][j] == 0) {
        excluded[i * numAccelerators + j] = 1
      }
    }
  }
}

dvar int numReplicas[vars];

minimize sum(v in vars) numReplicas[v] * costVector[v];
subject to {
  forall(i in servers) {
    sum(v in vars) numReplicas[v] * rateVector[i][v] >= arrivalRates[i];
  }
  sum(v in vars) numReplicas[v] * excluded[v] == 0;
  forall(v in vars) {
    numReplicas[v] >= 0;
  }
};

execute{
	writeln("numReplicas =" + numReplicas);
}



