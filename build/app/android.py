# app/android.py
import os
import subprocess

from app.build import Builder


class AndroidBuilder(Builder):

    def before_build(self):
        super().before_build()
        self.prepare_gomobile()

    def build(self):
        self.before_build()

        clean_files = ["libXray-sources.jar", "libXray.aar"]
        self.clean_lib_files(clean_files)
        os.chdir(self.lib_dir)
        
        # Set environment variables for 16KB page alignment (Android 15+ compatibility)
        env = os.environ.copy()
        
        # Enhanced CGO flags for 16KB page alignment
        cgo_cflags = "-O2 -g"
        cgo_ldflags = (
            "-Wl,-z,max-page-size=0x4000 "
            "-Wl,-z,common-page-size=0x4000 "
            "-Wl,-z,separate-loadable-segments"
        )
        
        # Set environment variables for gomobile
        env["CGO_CFLAGS"] = cgo_cflags
        env["CGO_CXXFLAGS"] = cgo_cflags
        env["CGO_LDFLAGS"] = cgo_ldflags
        
        # Enhanced linker flags for 16KB alignment
        ldflags = (
            "-s -w "
            "-extldflags \""
            "-Wl,-z,max-page-size=0x4000 "
            "-Wl,-z,common-page-size=0x4000 "
            "-Wl,-z,separate-loadable-segments"
            "\""
        )
        
        # Build with Android API 21+ and enhanced 16KB page alignment support
        print("Building with 16KB page alignment for Android 15+ compatibility...")
        ret = subprocess.run(
            [
                "gomobile", "bind", 
                "-target", "android", 
                "-androidapi", "21",
                "-ldflags", ldflags
            ],
            env=env
        )
        if ret.returncode != 0:
            raise Exception("build failed")
