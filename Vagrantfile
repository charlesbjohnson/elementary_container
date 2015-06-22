plugin_manager = 'vagrant-multiplug'
unless Vagrant.has_plugin?(plugin_manager)
  system("vagrant plugin install #{plugin_manager}")
  exec("vagrant #{ARGV.join(' ')}")
end

require 'active_support/core_ext/string/strip'
require 'pathname'

Vagrant.configure(2) do |config|
  config.plugin.add_dependency 'activesupport', '4.1.11'
  config.plugin.add_dependency 'vagrant-env', '0.0.2'

  config.env.enable

  config.vm.define 'builder' do |builder|
    builder.vm.box = 'phusion/ubuntu-14.04-amd64'

    go_path = Pathname.new(ENV['GOPATH'])
    go_project_path = Pathname.pwd
    vagrant_go_path = Pathname.new('/home/vagrant/go')
    vagrant_go_project_path = vagrant_go_path.join(go_project_path.relative_path_from(go_path))

    builder.vm.synced_folder '.', '/vagrant', disabled: true
    builder.vm.synced_folder '.', vagrant_go_project_path.to_s

    builder.vm.network 'forwarded_port', guest: ENV['APP_PORT'], host: 9001

    builder.vm.provision :shell, inline: GoSetup.dependencies
    builder.vm.provision :shell, inline: GoSetup.install

    builder.vm.provision :shell, privileged: false, inline: GoSetup.environment
    builder.vm.provision :shell, privileged: false, inline: ProjectSetup.dependencies
  end

  config.vm.define 'target' do |target|
    target.vm.box = 'phusion/ubuntu-14.04-amd64'

    target.vm.synced_folder '.', '/vagrant', disabled: true
    target.vm.synced_folder './bin', '/home/vagrant/bin'

    target.vm.network 'forwarded_port', guest: ENV['APP_PORT'], host: 9002
  end
end

module GoSetup
  extend self

  def dependencies
    <<-SCRIPT.strip_heredoc
      apt-get update -qq
      apt-get install -qq git mercurial
    SCRIPT
  end

  def install
    <<-SCRIPT.strip_heredoc
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
    <<-SCRIPT.strip_heredoc
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
    <<-SCRIPT.strip_heredoc
      go get -u github.com/constabulary/gb/...

      go get -u github.com/nsf/gocode
      go get -u github.com/k0kubun/pp
      go get -u golang.org/x/tools/cmd/godoc
      go get -u github.com/motemen/gore

      go get -u github.com/derekparker/delve/cmd/dlv
    SCRIPT
  end
end
