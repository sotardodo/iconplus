pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'sotar12'
        NAMESPACE = 'sotar'
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
    }

    stages {
        stage('1. Checkout Code') {
            steps {
                echo "📥 Mengambil kode terbaru dari repository..."
                checkout scm
            }
        }

        stage('2. Build Docker Images') {
            parallel {
                stage('Build Frontend') {
                    steps {
                        echo "🔨 Building Frontend Docker Image..."
                        dir('frontend') {
                            script {
                                appFrontend = docker.build("${env.DOCKER_REGISTRY}/frontend:latest")
                            }
                        }
                    }
                }
                stage('Build Laravel') {
                    steps {
                        echo "🔨 Building Laravel Docker Image..."
                        dir('laravel') {
                            script {
                                appLaravel = docker.build("${env.DOCKER_REGISTRY}/laravel:latest")
                            }
                        }
                    }
                }
                stage('Build Golang') {
                    steps {
                        echo "🔨 Building Golang Docker Image..."
                        dir('go') {
                            script {
                                appGolang = docker.build("${env.DOCKER_REGISTRY}/golang:latest")
                            }
                        }
                    }
                }
            }
        }

        stage('3. Push to Docker Hub') {
            steps {
                echo "🚀 Mengunggah Docker Images ke Registry..."
                script {
                    docker.withRegistry('https://registry.hub.docker.com', "${env.DOCKER_CREDENTIALS_ID}") {
                        appFrontend.push("latest")
                        appLaravel.push("latest")
                        appGolang.push("latest")
                    }
                }
            }
        }

        stage('4. Deploy to Kubernetes') {
            steps {
                echo "☸️ Menerapkan pembaruan ke Kubernetes Cluster..."
                script {
                    // Masuk ke folder k8s lalu jalankan script deploy.sh atau apply manual
                    dir('k8s') {
                        // Memastikan namespace dan konfigurasi dasar aman
                        sh 'kubectl apply -f 00-namespace.yaml'
                        
                        // Melakukan restart deployment agar pod mengambil image terbaru (Rolling Update)
                        sh 'kubectl rollout restart deployment/frontend -n ${NAMESPACE}'
                        sh 'kubectl rollout restart deployment/laravel -n ${NAMESPACE}'
                        sh 'kubectl rollout restart deployment/golang -n ${NAMESPACE}'
                        
                        echo "⏳ Menunggu proses rollout selesai..."
                        sh 'kubectl rollout status deployment/frontend -n ${NAMESPACE}'
                        sh 'kubectl rollout status deployment/laravel -n ${NAMESPACE}'
                        sh 'kubectl rollout status deployment/golang -n ${NAMESPACE}'
                    }
                }
            }
        }
    }

    post {
        success {
            echo "🎉 Pipeline CI/CD Berhasil! Seluruh layanan telah diperbarui di Kubernetes."
        }
        failure {
            echo "❌ Pipeline CI/CD Gagal! Silakan periksa log di atas."
        }
    }
}