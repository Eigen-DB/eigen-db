import random

import dagger
from dagger import dag, function, object_type


@object_type
class EigenDb:
    @function
    async def publish(self, source: dagger.Directory) -> str:
        """Publish the application container after building and testing it on-the-fly"""
        await self.unit_test(source)
        await self.integration_test(source)
        return await self.build(source).publish(
            f"ttl.sh/eigen_db_{random.randrange(10 ** 8)}"
        )

    @function
    def build(self, source: dagger.Directory) -> dagger.Container:
        """Build the application container"""
        build = (
            self.build_env(source)
            .with_exec(["go", "build", "-o", "./bin/eigen_db"])
            .with_exec(["cp", "-r", "./eigen", "./bin"])
            .directory("./bin")
        )
        return ( # dont run prod container as root
            dag.container()
            .from_("alpine:3.20")
            .with_exec(["mkdir", "/app"])
            .with_directory("/app", build)
            .with_exposed_port(8080)
            .with_exec(["/app/eigen_db"])
        )

    @function
    async def unit_test(self, source: dagger.Directory) -> int:
        """Return the exit code of running unit tests"""
        return await (
            self.build_env(source)
            .with_exec(["useradd", "tester"])
            .with_user("tester")
            .with_exec(["go", "test", "./...", "-count=1", "-v"])
            .exit_code()
        )
    
    @function
    async def integration_test(self, source: dagger.Directory) -> int:
        """Return the exit code of running integration tests"""

        # start the container
        return await (
            self.build_env(source)
            .with_env_variable("TEST_MODE", "1")
            .with_exec([
                "curl", 
                "https://github.com/ovh/venom/releases/download/v1.1.0/venom.linux-amd64", "-Lo", "./venom", 
                "&&", 
                "chmod", "+x", "./venom"
            ])
            .with_exec(["go", "run", "main.go", "&"])
            .with_exec(["./venom", "run", "integration_tests/", "--output-dir=integration_tests/logs"])
            .exit_code()
        )

    @function
    def build_env(self, source: dagger.Directory) -> dagger.Container:
        """Build a ready-to-use development environment"""
        return (
            dag.container()
            .from_("golang:1.20")
            .with_directory("/app", source)
            .with_workdir("/app")
            .with_exec(["go", "mod", "download"])
        )