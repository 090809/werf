module Dapp
  module Builder
    class Base
      include CommonHelper

      attr_reader :docker
      attr_reader :conf

      STAGES_DEPENDENCIES = {
          prepare: nil,
          infra_install: :prepare,
          sources_1: :infra_install,
          infra_setup: :sources_1,
          app_install: :infra_setup,
          sources_2: :app_install,
          app_setup: :sources_2,
          sources_3: :app_setup,
          sources_4: :sources_3
      }.freeze

      STAGES_DEPENDENCIES.each do |stage, dependence|
        define_method :"#{stage}_from" do
          send(:"#{dependence}_key") unless dependence.nil?
        end

        define_method :"#{stage}_image" do
          "dapp:#{send(:"#{stage}_key")}"
        end

        define_method :"#{stage}!" do
          build_stage!(from: send(:"#{stage}_from"), stage: stage)
        end
      end

      [:infra_install, :infra_setup].each do |stage|
        define_method :"#{stage}_key" do
          send("#{stage}_from")
        end
      end

      # TODO Describe stages sequence with
      # TODO   ordering data
      # TODO Generate stages related methods
      # TODO   from that data

      def initialize(docker, conf)
        @docker = docker
        @conf = conf
      end

      def run
        if prepare?
          prepare!
          infra_install!
          sources_1!
          infra_setup!
          app_install!
          app_setup!
        elsif infra_install?
          infra_install!
          sources_1!
          infra_setup!
          app_install!
          app_setup!
        elsif infra_setup?
          infra_setup!
          app_install!
          sources_2!
          app_setup!
        elsif app_install?
          app_install!
          sources_2!
          app_setup!
        elsif app_setup?
          app_setup!
          sources_3!
          sources_4!
        end
      end

      def build_stage!(from:, stage:)
        raise
      end


      def prepare?
        raise
      end

      def prepare
        # запуск shell-команд из conf
      end

      def prepare_from
        conf[:from]
      end

      def prepare_key
        # hash от shell-команд
      end


      def infra_install?
        raise
      end

      def infra_install
        raise
      end


      def infra_setup?
        raise
      end

      def infra_setup
        raise
      end


      def app_install?
        raise
      end

      def app_install
        raise
      end

      def app_install_key
        raise
      end


      def app_setup?
        raise
      end

      def app_setup
        raise
      end

      def app_setup_key
        raise
      end


      def sources_1
        raise
      end

      def sources_1_key
        raise
      end


      def sources_2
        raise
      end

      def sources_2_key
        raise
      end


      def sources_3
        raise
      end

      def sources_3_key
        raise
      end


      def sources_4
        raise
      end

      def sources_4_key
        raise
      end
    end
  end
end
