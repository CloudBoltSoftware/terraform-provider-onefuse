/*
    A pipeline to build the OneFuse Terraform Provider

    TODO:
    * 

    Jenkins Prerequisites:
    * 

    Agent Prerequisites:
    * Agent built with the go toolchain
    * Agent containing label 'go'
*/
pipeline {
    agent {
      node { label 'go' }
    }
    parameters {
        string(name: 'version', defaultValue: 'X.Y.Z', description: 'One Fuse Terraform Provider Version')
        string(name: 'bucket', defaultValue: 'cb-internal-builds', description: 'Bucket for uploading release artifacts.')
        string(name: 'bucket_root_path', defaultValue: '/OneFuse/Terraform/', description: 'Root path in bucket. "/" is main bucket as root.')
        string(name: 'release_date', defaultValue: 'YYYY-MM-DD', description: 'Release date of artifact.')
    }
    environment {
      TERRAFORM_PROVIDER_DIR = "terraform-provider-onefuse"
      TERRAFORM_PROVIDER_BIN_NAME = "terraform-provider-onefuse"
    }
    stages {
        stage('Build') {
            environment {
                  // Set the Go environment variables to be relative to the workspace directory
                  GOPATH = "${env.WORKSPACE}/go"
                  GOCACHE = "${env.WORKSPACE}/go/.cache"
            }
            steps {
                dir("${env.TERRAFORM_PROVIDER_DIR}") {
                    dir("linux") { } 
                    sh "GOOS=linux GOARCH=amd64 go build -o linux/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}"
                    writeFile(
                        file: "linux/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.sha256sum",
                        text: sh(
                            script: "sha256sum linux/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}",
                            returnStdout: true,
                        ).trim(),
                    )

                    dir("windows") { }
                    sh "GOOS=windows GOARCH=amd64 go build -o windows/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.exe"
                    writeFile(
                        file: "windows/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.sha256sum",
                        text: sh(
                            script: "sha256sum windows/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.exe",
                            returnStdout: true,
                        ).trim(),
                    )

                    dir("darwin") { }
                    sh "GOOS=darwin GOARCH=amd64 go build -o darwin/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}"
                    writeFile(
                        file: "darwin/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.sha256sum",
                        text: sh(
                            script: "sha256sum darwin/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}",
                            returnStdout: true,
                        ).trim(),
                    )
                }
            }
        }
        stage('Generate info.json') {
            environment {
                    CB_BUILD = "${env.GIT_COMMIT[0..9]}"

                    TERRAFORM_PROVIDER_BIN_FILE_PATH = "${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}"

                    SHA_256_CHECKSUM_LINUX = sh(
                        script: "cat ${env.TERRAFORM_PROVIDER_DIR}/linux/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.sha256sum | cut -d ' ' -f 1",
                        returnStdout: true,
                    ).trim()

                    SHA_256_CHECKSUM_DARWIN = sh(
                        script: "cat ${env.TERRAFORM_PROVIDER_DIR}/darwin/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.sha256sum | cut -d ' ' -f 1",
                        returnStdout: true,
                    ).trim()

                    SHA_256_CHECKSUM_WINDOWS = sh(
                        script: "cat ${env.TERRAFORM_PROVIDER_DIR}/windows/${env.TERRAFORM_PROVIDER_BIN_NAME}_v${params.version}.sha256sum | cut -d ' ' -f 1",
                        returnStdout: true,
                    ).trim()

            }
            steps {
                dir(WORKSPACE) {
                    writeFile(
                        file: "${env.TERRAFORM_PROVIDER_DIR}/info.json",
                        text: sh(
                            script: "./create_info.sh ${params.version} ${env.CB_BUILD} ${params.release_date} ${TERRAFORM_PROVIDER_BIN_FILE_PATH} ${env.SHA_256_CHECKSUM_LINUX} ${env.SHA_256_CHECKSUM_DARWIN} ${env.SHA_256_CHECKSUM_WINDOWS}",
                            returnStdout: true,
                        ).trim(),
                    )
                }
            }
        }
        stage("Upload release artifacts to S3") {
            steps {
                withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', accessKeyVariable: 'AWS_ACCESS_KEY_ID', credentialsId: 'AWS Jenkins User', secretKeyVariable: 'AWS_SECRET_ACCESS_KEY']]) {
                    sh script: "aws s3 sync ${env.TERRAFORM_PROVIDER_DIR} s3://${params.bucket}${bucket_root_path}${params.version} --exclude=* --include=linux/*${params.version}* --include=darwin/*${params.version}* --include=windows/*${params.version}* --include=info.json"
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: "${env.TERRAFORM_PROVIDER_DIR}/linux/**"
            archiveArtifacts artifacts: "${env.TERRAFORM_PROVIDER_DIR}/darwin/**"
            archiveArtifacts artifacts: "${env.TERRAFORM_PROVIDER_DIR}/windows/**"
            archiveArtifacts artifacts: "${env.TERRAFORM_PROVIDER_DIR}/info.json"
        }
    }

}
