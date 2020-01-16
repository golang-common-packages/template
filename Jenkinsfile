/* Only keep the 10 most recent builds and config schedule 10m. */
def projectProperties = [
    [$class: 'BuildDiscarderProperty',strategy: [$class: 'LogRotator', numToKeepStr: '10']],
	pipelineTriggers([pollSCM('*/10 * * * *')])
]

properties(projectProperties)

node {
  try {
    notifyBuild('STARTED')

    stage('Clone Backend golang-common-packages/template Repository'){
		  checkout scm
	  }

    stage('Build Docker Image for golang-common-packages/template') {
      sh 'docker build --rm -t golang-common-packages:latest .'
		}

    stage('Remove golang-common-packages/template Container') {
      sh 'docker-compose down'
		}    
	
    stage('Run Docker Image for golang-common-packages/template') {
      sh 'docker-compose up -d'
	  }
  } catch (e) {
      currentBuild.result = "FAILED"
  } finally {
      notifyBuild(currentBuild.result)
  }
}

def notifyBuild(String buildStatus = 'STARTED') {
  // build status of null means successful
  buildStatus =  buildStatus ?: 'SUCCESSFUL'

  // Default values
  def colorName = 'RED'
  def colorCode = '#FF0000'
  def subject = "${buildStatus}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'"
  def summary = "${subject} <${env.BUILD_URL}|Job URL> - <${env.BUILD_URL}/console|Console Output>"

  // Override default values based on build status
  if (buildStatus == 'STARTED') {
    color = 'YELLOW'
    colorCode = '#FFFF00'
  } else if (buildStatus == 'SUCCESSFUL') {
    color = 'GREEN'
    colorCode = '#00FF00'
  } else {
    color = 'RED'
    colorCode = '#FF0000'
  }

  // Send notifications
  slackSend(color: colorCode, message: summary)
}