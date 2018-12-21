ENV["LC_ALL"] = "en_US.UTF-8"

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = 'vitess'
  config.disksize.size = '50GB'

  config.vm.network "private_network", type: "dhcp"

  # vtctld
  config.vm.network "forwarded_port", guest: 8000, host: 8000 # http
  config.vm.network "forwarded_port", guest: 15000, host: 15000 # http
  config.vm.network "forwarded_port", guest: 15999, host: 15999 # grpc

  # vtgate
  config.vm.network "forwarded_port", guest: 15001, host: 15001 # http
  config.vm.network "forwarded_port", guest: 15991, host: 15991 # grpc
  config.vm.network "forwarded_port", guest: 15306, host: 15306 # mysql

  for i in 15000..17603
    config.vm.network :forwarded_port, guest: i, host: i
  end

  # Demo Appp
  config.vm.network "forwarded_port", guest: 8000, host: 8000 # http

  # If possible, use nfs, this gives a good boost to IO operations in the VM.
  # if you run into with nfs, just remove this from the synced folder

  config.vm.synced_folder ".", "/vagrant/src/vitess.io/vitess", type: "nfs"

  config.vm.provider :virtualbox do |vb|
    vb.name = "vitess"
    vb.customize ["modifyvm", :id, "--ioapic", "on"]
    vb.customize ["modifyvm", :id, "--cpuexecutioncap", "85"]
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.memory = 12888
    vb.cpus = 4
  end
  config.vm.provision "shell", path: "./vagrant-scripts/bootstrap_vm.sh"
end
