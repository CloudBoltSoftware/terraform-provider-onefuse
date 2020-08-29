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
        string(name: 'bucket', defaultValue: "internal-builds.cloudbolt.io", description: 'Bucket for uploading release artifacts.')
        string(name: 'bucket_root_path', defaultValue: '/OneFuse/Terraform/', description: 'Root path in bucket. "/" is main bucket as root.')
    }
    environment {
      VERSION = sh(
          script: "cat VERSION",
          returnStdout: true,
      ).trim()
      TAG = sh(returnStdout: true, script: "git tag --contains | head -1").trim()
      OUTPUT_BASEDIR = "release"
      OUTPUT_DIR = "${env.OUTPUT_BASEDIR}/terraform-provider-onefuse"
      TERRAFORM_PROVIDER_BIN_NAME = "terraform-provider-onefuse_v${env.VERSION}"
      DATE= sh(
	  script: "date +\"%m-%d-%y\"",
	  returnStdout: true,
      ).trim()
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
                            script: "./scripts/create_info.sh ${env.VERSION} ${env.CB_BUILD} ${env.DATE} ${TERRAFORM_PROVIDER_BIN_FILE_PATH} ${env.SHA_256_CHECKSUM_LINUX} ${env.SHA_256_CHECKSUM_DARWIN} ${env.SHA_256_CHECKSUM_WINDOWS}",
                            returnStdout: true,
                        ).trim(),
                    )
                }
            }
        }
        stage("Upload release artifacts to S3") {
            steps {
                withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', accessKeyVariable: 'AWS_ACCESS_KEY_ID', credentialsId: 'AWS Jenkins User', secretKeyVariable: 'AWS_SECRET_ACCESS_KEY']]) {
                    sh script: "aws s3 sync ${env.OUTPUT_DIR} s3://${params.bucket}${params.bucket_root_path}${env.VERSION} --exclude=* --include=linux/*${env.VERSION}* --include=darwin/*${env.VERSION}* --include=windows/*${env.VERSION}* --include=info.json"
                }
            }
        }
	stage('Send slack message') {
	    steps {
		slackSend(
		    channel: '#automation-testing-ground',
		    color: 'good',
		    blocks:[
		    [
			'type': 'header',
			'text': [
			    'type': 'plain_text',
			    'text': "OneFuse-Terraform-Provider ${GIT_BRANCH}-${BUILD_NUMBER} is here :meow_party:",
			    'emoji': true
			]
		    ],
		    [
			'type': 'section',
			'text': [
			    'type': 'mrkdwn',
			    'text': "s3://${params.bucket}${params.bucket_root_path}${env.VERSION}/${env.VMOAPP_NAME}"
			]
		    ],
		    [
			'type': 'divider'
		    ],
		    [
			'type': 'context',
			'elements': [
			    [
				'type': 'image',
				'image_url': 'https://pbs.twimg.com/profile_images/625633822235693056/lNGUneLX_400x400.jpg',
				'alt_text': 'cute cat'
			    ],
			    [
				'type': 'mrkdwn',
				'text': 'This is an internal only release candidate'
			    ]
			]
			]
		    ]
        )
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
