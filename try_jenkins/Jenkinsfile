pipeline {
    agent any
    environment {
        DOCKER_IMAGE = "main-server"
        DOCKER_TAG = "latest"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build'){
            steps {
                script {
                    docker.build(DOCKER_IMAGE)
                }
            }
        }

        stage('Test'){
            script{
                docker.image(DOCKER_IMAGE).inside{
                    sh 'go test ./...'
                }
            }
        }

        stage('Deploy'){
            steps{
                echo 'Deploy step (implement if needed)'
            }
        }
    }
}