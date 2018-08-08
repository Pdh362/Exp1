# Exp1

### Installation
No libraries are required to run the final executable. 

### Usage

The executable can be run in one of two modes :-

* Master - In this mode the executable awaits information from watchers, and collates the file listings into one master list.
* Watcher - In this mode a single file directory is monitored for any changes. Updates are forwarded via http POST command to the master.

Only one master should ever be active at a time - multiple watchers can be active. 

### Command line parameters
* _-help_    Print command-line information.
* _-mode_    Defines the mode. Valid options are 'watcher', 'master'.
* _-refresh_  Specifies how often a watcher node will update the master. 
* _-wport_     Defines which port a watcher runs on. **Must be unique per watcher**.
* _-mport_    Defines which port a master presents the collated information on.  

### Example

First, run your master : `exp1 -mode=master -mport=90` 

Then, as many watchers as you like :
* `exp1 -mode=watcher -watch=c:\ -wport=8000`
* `exp1 -mode=watcher -watch=d:\ -wport=8001`

Finally, go to `http://localhost:90` to view the results.

In each mode, the executable will block: that is, the code will run eternally until an error occurs, or the executable is halted.

## Improvements
* The code is hardwired to `localhost`: configuration option to change this would be useful.
* As commented in the code, there are a couple of areas of golang code that I feel could/should be improved.
* Performance at high watcher counts is untested :blush:
