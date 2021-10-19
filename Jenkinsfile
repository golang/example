//currentBuild.displayName = currentBuild.number+"#"+" testcase/coloraize/log_parser/build_name:"+currentBuild.number
//currentBuild.description = "test desc" 
currentBuild.description = ""
if (currentBuild.rawBuild.getCause(Cause).properties.upstreamRun != null) currentBuild.description = "Started by: ${currentBuild.rawBuild.getCause(Cause).properties.upstreamRun}<br>"
else if (currentBuild.rawBuild.getCause(Cause).properties.userName != null) currentBuild.description = "Started by: ${currentBuild.rawBuild.getCause(Cause).properties.userName}<br>"

node("test_node") {

    stage('git clone') { // for display purposes
	checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: 'ssh-hetzner', url: 'git@github.com:h34dl355/example.git']]])
	sh 'ls -la'
	sh 'id'
    	echo "Git clone done!"
    stage('test job') {
        echo "test"
    }
    }
    stage('build') {
	sh 'docker build . -t gendalf'
	echo "Build done!"
    }
    stage('docker run') {
        sh 'docker run gendalf'
	echo "End of pipeline"
    }
}
