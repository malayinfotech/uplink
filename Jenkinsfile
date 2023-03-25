pipeline {
    agent {
        docker {
            label 'main'
            image 'storxlabs/ci:latest'
            alwaysPull true
            args '-u root:root --cap-add SYS_PTRACE -v "/tmp/gomod":/go/pkg/mod'
        }
    }
    options {
          timeout(time: 26, unit: 'MINUTES')
    }
    environment {
        NPM_CONFIG_CACHE = '/tmp/npm/cache'
        COCKROACH_MEMPROF_INTERVAL=0
    }
    stages {
        stage('Build') {
            steps {
                checkout scm

                sh 'mkdir -p .build'
                // make a backup of the mod file in case, for later linting
                sh 'cp go.mod .build/go.mod.orig'
                sh 'cp testsuite/go.mod .build/testsuite.go.mod.orig'

                sh 'service postgresql start'

                dir(".build") {
                    sh 'cockroach start-single-node --insecure --store=\'/tmp/crdb\' --listen-addr=localhost:26257 --http-addr=localhost:8080 --cache 512MiB --max-sql-memory 512MiB --background'
                }
            }
        }

        stage('Verification') {
            parallel {
                stage('Lint') {
                    steps {
                        sh 'check-copyright'
                        sh 'check-large-files'
                        sh 'check-imports ./...'
                        sh 'check-peer-constraints'
                        sh 'storx-protobuf --protoc=$HOME/protoc/bin/protoc lint'
                        sh 'storx-protobuf --protoc=$HOME/protoc/bin/protoc check-lock'
                        sh 'check-atomic-align ./...'
                        sh 'check-monkit ./...'
                        sh 'check-errs ./...'
                        sh './scripts/check-dependencies.sh'
                        sh 'staticcheck ./...'
                        sh 'golangci-lint --config /go/ci/.golangci.yml -j=2 run'
                        sh 'go-licenses check ./...'
                        sh './scripts/check-libuplink-size.sh'
                        sh 'check-mod-tidy -mod .build/go.mod.orig'

                        dir("testsuite") {
                            sh 'check-imports ./...'
                            sh 'check-atomic-align ./...'
                            sh 'check-monkit ./...'
                            sh 'check-errs ./...'
                            sh 'staticcheck ./...'
                            sh 'golangci-lint --config /go/ci/.golangci.yml -j=2 run'
                            sh 'check-mod-tidy -mod ../.build/testsuite.go.mod.orig'
                        }
                    }
                }

                stage('Tests') {
                    environment {
                        COVERFLAGS = "${ env.BRANCH_NAME == 'main' ? '-coverprofile=.build/coverprofile -coverpkg=./...' : ''}"
                    }
                    steps {
                        sh 'go vet ./...'
                        sh 'go test -parallel 4 -p 6 -vet=off $COVERFLAGS -timeout 20m -json -race ./... > .build/tests.json'
                    }

                    post {
                        always {
                            sh script: 'cat .build/tests.json | xunit -out .build/tests.xml', returnStatus: true
                            sh script: 'cat .build/tests.json | tparse -all -top -slow 100', returnStatus: true
                            archiveArtifacts artifacts: '.build/tests.json'
                            junit '.build/tests.xml'

                            script {
                                if(fileExists(".build/coverprofile")){
                                    sh script: 'filter-cover-profile < .build/coverprofile > .build/clean.coverprofile', returnStatus: true
                                    sh script: 'gocov convert .build/clean.coverprofile > .build/cover.json', returnStatus: true
                                    sh script: 'gocov-xml  < .build/cover.json > .build/cobertura.xml', returnStatus: true
                                    cobertura coberturaReportFile: '.build/cobertura.xml'
                                }
                            }
                        }
                    }
                }

                stage('Testsuite') {
                    environment {
                        STORX_TEST_COCKROACH = 'cockroach://root@localhost:26257/testcockroach?sslmode=disable'
                        STORX_TEST_POSTGRES = 'postgres://postgres@localhost/teststorx?sslmode=disable'
                        STORX_TEST_COCKROACH_NODROP = 'true'
                        STORX_TEST_LOG_LEVEL = 'info'
                        COVERFLAGS = "${ env.BRANCH_NAME == 'main' ? '-coverprofile=../.build/testsuite_coverprofile -coverpkg=uplink/...' : ''}"
                    }
                    steps {
                        sh 'cockroach sql --insecure --host=localhost:26257 -e \'create database testcockroach;\''
                        sh 'psql -U postgres -c \'create database teststorx;\''
                        dir('testsuite'){
                            sh 'go vet ./...'
                            sh 'go test -parallel 4 -p 6 -vet=off $COVERFLAGS -timeout 20m -json -race ./... > ../.build/testsuite.json'
                        }
                    }

                    post {
                        always {
                            dir('testsuite'){
                                sh script: 'cat ../.build/testsuite.json | xunit -out ../.build/testsuite.xml', returnStatus: true
                            }
                            sh script: 'cat .build/testsuite.json | tparse -all -top -slow 100', returnStatus: true
                            archiveArtifacts artifacts: '.build/testsuite.json'
                            junit '.build/testsuite.xml'

                            script {
                                if(fileExists(".build/testsuite_coverprofile")){
                                    sh script: 'filter-cover-profile < .build/testsuite_coverprofile > .build/clean.testsuite_coverprofile', returnStatus: true
                                    sh script: 'gocov convert .build/clean.testsuite_coverprofile > .build/testsuite_cover.json', returnStatus: true
                                    sh script: 'gocov-xml  < .build/testsuite_cover.json > .build/testsuite_cobertura.xml', returnStatus: true
                                    cobertura coberturaReportFile: '.build/testsuite_cobertura.xml'
                                }
                            }
                        }
                    }
                }

                stage('Integration [storx/storx]') {
                    environment {
                        STORX_TEST_POSTGRES = 'postgres://postgres@localhost/teststorx2?sslmode=disable'
                        STORX_TEST_COCKROACH = 'omit'
                        // TODO add 'omit' for metabase STORX_TEST_DATABASES
                        STORX_TEST_DATABASES = 'pg|pgx|postgres://postgres@localhost/testmetabase?sslmode=disable'
                    }
                    steps {
                        sh 'psql -U postgres -c \'create database teststorx2;\''
                        sh 'psql -U postgres -c \'create database testmetabase;\''
                        dir('testsuite'){
                            sh 'cp go.mod go-temp.mod'
                            sh 'go vet -modfile go-temp.mod -mod=mod storx/...'
                            sh 'go test -modfile go-temp.mod -mod=mod -parallel 4 -p 6 -vet=off -timeout 20m -json storx/... > ../.build/testsuite-storx.json'
                        }
                    }

                    post {
                        always {
                            dir('testsuite'){
                                sh 'cat ../.build/testsuite-storx.json | xunit -out ../.build/testsuite-storx.xml'
                            }
                            sh script: 'cat .build/testsuite-storx.json | tparse -all -top -slow 100', returnStatus: true
                            archiveArtifacts artifacts: '.build/testsuite-storx.json'
                            junit '.build/testsuite-storx.xml'
                        }
                    }
                }

                stage('Integration [rclone]') {
                    environment {
                        STORX_SIM_POSTGRES = 'postgres://postgres@localhost/teststorx3?sslmode=disable'
                    }
                    steps {
                        sh 'psql -U postgres -c \'create database teststorx3;\''
                        echo 'Testing against satellite'
                        sh './testsuite/scripts/test-sim.sh'
                        sh 'psql -U postgres -c \'drop database teststorx3;\''
                    }
                    post {
                        always {
                            zip zipFile: 'rclone-integration-tests.zip', archive: true, dir: '.build/rclone-integration-tests'
                        }
                    }
                }

                stage('Go Compatibility') {
                    steps {
                        sh 'check-cross-compile -compiler "go,go.min" uplink/...'
                    }
                }
            }
        }
    }

    post {
        always {
            sh "chmod -R 777 ." // ensure Jenkins agent can delete the working directory
            deleteDir()
        }
    }
}
