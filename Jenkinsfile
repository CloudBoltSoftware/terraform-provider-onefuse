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
        string(name: 'version', defaultValue: 'X.Y.Z', description: 'IGNORED. Version is set in the VERSION file of the repository.')
        string(name: 'bucket', defaultValue: 'cb-internal-builds', description: 'Bucket for uploading release artifacts.')
        string(name: 'bucket_root_path', defaultValue: '/OneFuse/Terraform/', description: 'Root path in bucket. "/" is main bucket as root.')
        string(name: 'release_date', defaultValue: 'YYYY-MM-DD', description: 'Release date of artifact.')
    }
    environment {
      VERSION = sh(
          script: "cat VERSION",
          returnStdout: true,
      ).trim()
      OUTPUT_BASEDIR = "release"
      OUTPUT_DIR = "${env.OUTPUT_BASEDIR}/terraform-provider-onefuse"
      TERRAFORM_PROVIDER_BIN_NAME = "terraform-provider-onefuse_v${env.VERSION}"
    }
    stages {
        stage('Build') {
            environment {
                  // Set the Go environment variables to be relative to the workspace directory
                  GOPATH = "${env.WORKSPACE}/go"
                  GOCACHE = "${env.WORKSPACE}/go/.cache"
            }
            steps {
                sh "make release"
            }
        }
        stage('Generate info.json') {
            environment {
                    CB_BUILD = "${env.GIT_COMMIT[0..9]}"

                    TERRAFORM_PROVIDER_BIN_FILE_PATH = "${env.TERRAFORM_PROVIDER_BIN_NAME}"

                    SHA_256_CHECKSUM_LINUX = sh(
                        script: "cat ${env.OUTPUT_DIR}/linux/${env.TERRAFORM_PROVIDER_BIN_NAME}.sha256sum | cut -d ' ' -f 1",
                        returnStdout: true,
                    ).trim()

                    SHA_256_CHECKSUM_DARWIN = sh(
                        script: "cat ${env.OUTPUT_DIR}/darwin/${env.TERRAFORM_PROVIDER_BIN_NAME}.sha256sum | cut -d ' ' -f 1",
                        returnStdout: true,
                    ).trim()

                    SHA_256_CHECKSUM_WINDOWS = sh(
                        script: "cat ${env.OUTPUT_DIR}/windows/${env.TERRAFORM_PROVIDER_BIN_NAME}.sha256sum | cut -d ' ' -f 1",
                        returnStdout: true,
                    ).trim()

            }
            steps {
                dir(WORKSPACE) {
                    writeFile(
                        file: "${env.OUTPUT_DIR}/info.json",
                        text: sh(
                            script: "./scripts/create_info.sh ${env.VERSION} ${env.CB_BUILD} ${params.release_date} ${TERRAFORM_PROVIDER_BIN_FILE_PATH} ${env.SHA_256_CHECKSUM_LINUX} ${env.SHA_256_CHECKSUM_DARWIN} ${env.SHA_256_CHECKSUM_WINDOWS}",
                            returnStdout: true,
                        ).trim(),
                    )
                }
            }
        }
        stage("Upload release artifacts to S3") {
            steps {
                withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', accessKeyVariable: 'AWS_ACCESS_KEY_ID', credentialsId: 'AWS Jenkins User', secretKeyVariable: 'AWS_SECRET_ACCESS_KEY']]) {
                    sh script: "aws s3 sync ${env.OUTPUT_DIR} s3://${params.bucket}${bucket_root_path}${env.VERSION} --exclude=* --include=linux/*${env.VERSION}* --include=darwin/*${env.VERSION}* --include=windows/*${env.VERSION}* --include=info.json"
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: "${env.OUTPUT_DIR}/linux/**"
            archiveArtifacts artifacts: "${env.OUTPUT_DIR}/darwin/**"
            archiveArtifacts artifacts: "${env.OUTPUT_DIR}/windows/**"
            archiveArtifacts artifacts: "${env.OUTPUT_DIR}/info.json"
        }
    }

}
