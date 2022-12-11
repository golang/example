@Library('shared-pipeline@feature/rewrite') _

  node("slave-vm") {
    timestamps {
      stage('Checkout project') {
        println "Executing example pipeline.."
        cleanWs()
        checkoutProject()
      }
      stage('Test stage') {
        dockerBuildPush()
        sh 'docker images'
      }
      println "Ended example pipeline."
    }
  }