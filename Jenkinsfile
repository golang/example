node("test_node") {

    stage('git clone') { // for display purposes
	checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: 'ssh-hetzner', url: 'git@github.com:h34dl355/example.git']]])
	sh 'ls -la'
	sh 'id'
    	echo "Git clone done!"
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
