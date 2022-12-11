@Library('shared-pipeline@feature/rewrite') _

  node(jenkinsAgent) {
    timestamps {
      stage('Checkout project') {
        println "Executing example pipeline.."
        cleanWs()
        checkoutProject()
      }
      stage('Test stage') {
        sh 'ls -la'
      }
    }
  }