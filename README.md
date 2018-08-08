# Exp1

### Installation
No libraries are required to run the final executable.

### Command line paramaters
* _-help_    Print command-line information
* _-mode_    Defines the mode. Valid options are 'watcher', 'master'
* _-refresh_  Specifies how often a watcher node will update the master. 
* _-wport_     Defines which port a watcher speaks to the master. **Must be unique per watcher**
* _-mport_    Defines which port a master presents the collated information to.  

### Usage

The executable can be run in one of two modes :-

* Master - In this mode the executable awaits information from watchers, and collates the file listings into one master list.
* Watcher - In this mode a single file directory is monitored for any changes. Updates are forwarded, via a http POST command, to the master.

Only one master should ever be active at a time - multiple watchers can be active. 