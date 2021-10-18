node("test_node") {

    stage('git clone') { // for display purposes
	checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: 'ssh-hetzner', url: 'git@github.com:h34dl355/example.git']]])
	sh 'ls -la'
    	echo "Git clone done!"
    }
    stage('build') {
	sh 'docker build . -t gendalf'
	echo "Build done!"
    }
    stage('docker run') {
        docker run gendalf
	echo "End of pipeline"
    }
}
