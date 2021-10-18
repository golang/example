node("test_node") {

    stage('first') { // for display purposes
	checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: 'ssh-hetzner', url: 'git@github.com:h34dl355/example.git']]])
	sh 'ls -la'
    	echo "Git clone done!"
    }
    stage('second') {
	sh 'docker run -v ~/root/example/hello/:/go/src/ golang:latest go build'
    }
}
