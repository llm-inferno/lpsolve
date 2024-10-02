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

int maxNumReplicas[servers][accelerators];

float costVector[vars];
int assignVector[servers][vars];
int excluded[vars];

int numReplicas[servers][accelerators];

execute {
  for(var i in servers) {
    for(var j in accelerators) {
      if (ratePerReplica[i][j] > 0) {
        maxNumReplicas[i][j] = Opl.ceil(arrivalRates[i] / ratePerReplica[i][j]);
      } else {        
        excluded[i * numAccelerators + j] = 1;
      }
    }
  }
}


execute {
  for(var i in servers) {
    for(var j in accelerators) {
      costVector[i * numAccelerators + j] = numInstancesPerReplica[i][j] * instanceCost[j] * maxNumReplicas[i][j];
      assignVector[i][i * numAccelerators + j] = 1;
    }
  }
}

dvar boolean assignment[vars];

minimize sum(v in vars) assignment[v] * costVector[v];
subject to {
  forall(i in servers) {
    sum(v in vars) assignment[v] * assignVector[i][v] == 1;
  }
  sum(v in vars) assignment[v] * excluded[v] == 0;
};

execute{
  for(var i in servers) {
    for(var j in accelerators) {
      numReplicas[i][j] = assignment[i * numAccelerators + j] * maxNumReplicas[i][j];
    }
  }  
  writeln("numReplicas =" + numReplicas);
}
