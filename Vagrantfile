require 'pathname'

Vagrant.configure(2) do |config|
  config.vm.define 'builder' do |builder|
    builder.vm.box = 'phusion/ubuntu-14.04-amd64'

    go_path = Pathname.new(ENV['GOPATH'])
    go_project_path = Pathname.pwd
    vagrant_go_path = Pathname.new('/home/vagrant/go')
    vagrant_go_project_path = vagrant_go_path.join(go_project_path.relative_path_from(go_path))

    builder.vm.synced_folder '.', '/vagrant', disabled: true
    builder.vm.synced_folder '.', vagrant_go_project_path.to_s

    builder.vm.provision :shell, inline: GoSetup.dependencies
    builder.vm.provision :shell, inline: GoSetup.install

    builder.vm.provision :shell, privileged: false, inline: GoSetup.environment
    builder.vm.provision :shell, privileged: false, inline: ProjectSetup.dependencies
  end

  config.vm.define 'target' do |target|
    target.vm.box = 'phusion/ubuntu-14.04-amd64'
    target.vm.synced_folder '.', '/vagrant', disabled: true
    target.vm.synced_folder './bin', '/home/vagrant/bin'
  end
end

def strip_heredoc(s)
  indent = s.scan(/^[ \t]*(?=\S)/).min.size
  s.gsub(/^[ \t]{#{indent}}/, '')
end

module GoSetup
  extend self

  def dependencies
    strip_heredoc(<<-SCRIPT)
      apt-get update -qq
      apt-get install -qq git mercurial
    SCRIPT
  end

  def install
    strip_heredoc(<<-SCRIPT)
      if [[ ! -d /usr/local/go ]]; then
        mkdir --parents /tmp
        curl --silent --location --output /tmp/go.tar.gz https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
        tar --extract --gunzip --file /tmp/go.tar.gz --directory /usr/local
        chown --recursive vagrant /usr/local/go
        rm /tmp/go.tar.gz
      fi
    SCRIPT
  end

  def environment
    strip_heredoc(<<-SCRIPT)
      GOROOT=/usr/local/go
      GOPATH=~/go

      mkdir --parents $GOPATH

      cat << EOF > ~/.profile.go
      export GOROOT=$GOROOT
      export GOPATH=$GOPATH
      export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
      EOF

      sed --in-place "/source ~\\/.profile.go/d" ~/.profile
      echo "source ~/.profile.go" >> ~/.profile

      unset GOROOT
      unset GOPATH
    SCRIPT
  end
end

module ProjectSetup
  extend self

  def dependencies
    strip_heredoc(<<-SCRIPT)
      go get -u github.com/constabulary/gb/...

      # TODO add gore
      # go get -u github.com/nsf/gocode
      # go get -u github.com/k0kubun/pp
      # go get -u golang.org/x/tools/cmd/godoc
      # go get -u github.com/motemen/gore

      go get -u github.com/derekparker/delve/cmd/dlv
    SCRIPT
  end
end
