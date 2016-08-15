module Dapp
  # Application
  class Application
    # Path
    module Path
      def home_path(*path)
        make_path(config._home_path, *path).expand_path
      end

      def tmp_path(*path)
        make_path(@tmp_path, *path).expand_path.tap { |p| p.parent.mkpath }
      end

      def build_path(*path)
        make_path(@build_path, home_path.basename, *path).expand_path.tap { |p| p.parent.mkpath }
      end

      def lock_path(*path)
        make_path(@lock_path, *path).expand_path
      end

      def container_dapp_path(*path)
        make_path('/.dapp', *path)
      end

      def container_tmp_path(*path)
        container_dapp_path('tmp', *path)
      end

      private

      def make_path(base, *path)
        path.compact.map(&:to_s).inject(Pathname.new(base), &:+)
      end
    end # Path
  end # Application
end # Dapp
