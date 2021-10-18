node("test_node") {

    stage('first') { // for display purposes
	checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: 'ssh-hetzner', url: 'git@github.com:h34dl355/example.git']]])
    	echo "Its work"
	sh 'ls -la'
    }
}
