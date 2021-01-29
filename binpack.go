package main

import (
	"archive/zip"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type BinpackFile struct {
	File     string
	MimeType string
	Content  string
}

var (
	binpackReadFiles *bool
	BinpackFiles     = make(map[string]*BinpackFile)

	Binpack_Icon_favicon16x16Png = &BinpackFile{
		File:     "icon/favicon-16x16.png",
		MimeType: "image/png",
		Content:  "89504e470d0a1a0a0000000d49484452000000100000001008060000001ff3ff610000017849444154384fad934b2844611480bf7fc61d06330d45364a344b21b351cac26362cfc26614839d0d3b5928d99185685066a13c16928528366c282c4899918812258f1ba3cbccd5bd5ed7848c7176e7f59d73feff1c41942457fb8b12226ab38a2807725edd470275f9c9247cf70b9e6d638a78576aa7adb6dbd000d0047cd83f17508151d96e6d63a62ea4b95e025f929780d2e88ebed1d764bbb54a83e8009bdbef03bcbf4c7e0b1b91173dcd429bd91c61f387b6bfe3aa6113c5c2ee1e1f5211ad3156d7c305eab0b0b9fd01c0a919667b2a58db3967efe89afcdc344e2fee58583fa5a1da89d9247850c25824137d53bb6ff5821a400124cdd2dd584cd7d8263d5e1737770a998e2446e6f7a929c9c69e22b11dbc642b70a9835fe5f14b406f8b8b2b596162e9804e4f1181931be6568f71d82cd457e6d13eb8f109f03e82afa3540f5ed93aa3d099cee4f22165055964a45971a45a780a474894ccf44f1b4688fb11e3fec6b817e95f56d900f9e33119d630d6737e06bfd4a7ef4dee3b410000000049454e44ae426082",
	}
	Binpack_Icon_favicon32x32Ico = &BinpackFile{
		File:     "icon/favicon-32x32.ico",
		MimeType: "image/x-icon",
		Content: "zip:" +
			"504b0304140008080800000000000000000000000000000000001600000069636f6e5c66617669636f6e2d33327833322e69636f626060646064505000d1120c2b781818c418181834181818141818181c1820e2e4014e16460b49d6246d8e464baea94e3c0bdc7917b8f34e75e269b4e44ad2e6b09464e564612464067620c1cd94acc331db8577a13b3e34db8537598743829b89907908c0c6cc18a1ce3ecf8d80c9c8689e1b6f843a3b1b3361bf887331b559731334102b6ab3e616e7c2e711453ee6a94e3c04cdc183a63af128f23163355c9c8b8942c3e15660fa828d9991ec60c1446dd6dc687111a1ce4e50174928429d1d6eb804371349a9851834cf8d179e6893753808aa270325eb70407228c14c441e9aedc2cbc9c26829c94a5025d90852b6a009dedbb7eac1a175f70fac393db362a99f185c7079882c84bdabc2f7ceae250bdd7997f8895d5cd6f9f6f68577f72e5fdb306d5980249a5190820b4df0efaf1f3b4bbdb6e639dcdeb1f0d5d5138bbd852082aba3d416baf32e0b94faf2f2d1cfcfef17baf3dedeb9e8defe552b42e49705493f38b4eefac6196846414a454cf321462df4e07b75fdd4a1f64464c11b9b66beba761262fec969a56b62b520ba4ecfaa7a7874339a51902217a7f9eebc57564fbcba6e0a5c707b91dbf7f72fb7153843cc87b86147b1fb8d4d337f7c" +
			"78b3b3dc07cd2848798ec7fccb2b7a21befefbebc7ba04dd0f8f6e1cee4c5997a80f317f75b4fa9b5be7deddbd746e7ec3ea68f5851846412a0b3ce6df3fb0e6d48c7288e0ed9d8b9f9cdab5d09d176efeb5f5d3eeec5ebad0836fa13befda789ded45ae98e183357e21e6ef2cf3fef6eee5aa481588e0cf2f1fd6c669239b7f7945efbd7dab1679f02ff2e4bfbb77c5955513d08c82547398e67f7c7cfbd3b37bafaf9fde9ae700173c39b504c2869bbf3c48fae1918d5f5e3dfef2f2f1bdfdabb0a64f0b5ae62f48054dd3f281d6e51b1dca675ad72f74a81f695dbf43004ddb271040d3f61504d0b47d0807b46bdf220348058dab7d6e4141fb9c5e00100000ffff504b0708458ec22fc3020000be0c0000504b0102140014000808080000000000458ec22fc3020000be0c000016000000000000000000000000000000000069636f6e5c66617669636f6e2d33327833322e69636f504b0506000000000100010044000000070300000000",
	}
	Binpack_Icon_favicon32x32Png = &BinpackFile{
		File:     "icon/favicon-32x32.png",
		MimeType: "image/png",
		Content: "89504e470d0a1a0a0000000d4948445200000020000000200806000000737a7af4000003e7494441545847c5977d4cd4751cc7dfdfdfef3c38e0b8094d1e26c99ca3166589e1034c22a8ce603d0cda0aca044df1613e4c24738a16b4351dac65f3ec56d0812bda6044292c1626f63088c09c4aa14b4bd9149053e988e3e0f87dddf77bdc0dee018e38f7fbfdf3db6fdfcfc3ebfbf97cbe9fdff743e0e533f719bd664c54a5532aa5508225046401809071f5db14f41aa1f89d0838255a2d0d779af206bc314da61352a755c4600c7b40910502d574f27c9dc20c822a8838646a587b792a1d8f00912fe8034c237e4500d90e608e578e5d8546017a44adb41cb871226fc89d0db700c1dae38b28a43a00b1ffd3b1b35a2781f0f2bf8d6bfe725e7001d03c5bf9a424d0ef0084fac8b9dd8c5190c8ea81efdf6c9f687712c0f8ce5bef8373070481b06262241c00b69cfbb7f930ec9e02d8a9560e2fb3d7840340ad35940024dfc761f7608e969a1a7376b3450e307ed42ecea2da67ca3d0a118fb2236a03d05694015837532bb3942f3735ae5d4f5887b30afe37bd6e32b3f4ea50a7302ba4e10812bcda904529f9d2577667628780bc4ed4cf193e0321eb9d156f9d7803ddbd83104502e3800505c7dad071a99f8bb1b5fd9fb643ff6d17ff5e18a946f391743cf8ca57080e5442b72b017131a1" +
			"ac25a3e58f3eecf8a81583e65157364acb4890d6d04140e2dc012ccea9c54de310921e0f47f93bab90b8f5247a6f9b39c0f08815895b4ee27aefe024807773e310352f101b0effcc4b5cbf3b113d46330acb3a5c0028e859a2d656b06db9743de6c40ec0343fdf9b84dfba6e41f7f59f1ca0faf45584cd5521b3f0d42480439be2a109526273e92fa01458101e8488d000b476f6b9cb8e910148f6e33851c21980ed2c50a54081ae8d03c46ffc06d5452928a9bac0c1ec29887c200055079f86421450fbe33fa86aba821bfd6eff43cc1df51ae0fdb796f27ad8ab6fe7002c3ad1e141385e988cecf74ea3a63895d7006f2e0458191b868ca7a29199148d9c0fcee0ccb91e7711e000d3a680196cfaf0791cabeb424df3df0e00561f255b97e38945218889d270804d2f3d8cfa966e74f7fdc71d6ecb7c042b63e721bba8d97d0aa62ac28c7d4d18b258919b168387a234487bbb11a3566912803a600e5a3f7911eccd008eee4a80bf5244c1d1365825091fef4c40df1d334f9df3632b420fc7b06ccf2a2814022489a2fd523fcaeb2fc36cb1721b6c2d5ff72bee9a46f877f29208bc9ab290179eca4f44fe6b8f411b3f9fa7eca7f33d28369cf37c0c83b595d914f48b9934105fc91242b3e56fc56c37b2fe8c3880ede62bdfefd81605192f240c40f62b198390f5526a3f5ab25ecbed10b20e" +
			"2676085947b3899d4eb6e1d4b9ddf2f15ce1974625a44e399e13e10771cc5cefed787e0febb2dd10ee52c6770000000049454e44ae426082",
	}
)

func init() {
	binpackReadFiles = flag.Bool("binpack.readfiles", !common.IsRunningAsExecutable(), "Read resource from filesystem")

	BinpackFiles["icon/favicon-16x16.png"] = Binpack_Icon_favicon16x16Png
	BinpackFiles["icon/favicon-32x32.ico"] = Binpack_Icon_favicon32x32Ico
	BinpackFiles["icon/favicon-32x32.png"] = Binpack_Icon_favicon32x32Png

	common.RegisterResourceLoader(func(name string) []byte {
		l := make([]string, 0)

		for k := range BinpackFiles {
			if filepath.Base(k) == filepath.Base(name) {
				l = append(l, k)
			}
		}

		if len(l) == 0 {
			return nil
		}

		if len(l) > 1 {
			common.Error(fmt.Errorf("multiple resources with name %s found: %v", name, l))
		}

		b := BinpackFiles[l[0]]

		ba, err := b.Unpack()
		if common.Error(err) {
			return nil
		}

		return ba
	})
}

func (this *BinpackFile) Unpack() ([]byte, error) {
	if *binpackReadFiles {
		filename := ""

		b := common.FileExists(filepath.Base(this.File))
		if b {
			filename = filepath.Base(this.File)
		} else {
			b = common.FileExists(this.File)
			if b {
				filename = this.File
			}
		}

		if b {
			common.Warn("Read resource from filesystem: %s", filename)

			ba, err := ioutil.ReadFile(filename)

			if err == nil {
				return ba, nil
			}
		}
	}

	if strings.HasPrefix(this.Content, "zip:") {
		ba, err := hex.DecodeString(this.Content[4:])
		if common.Error(err) {
			return ba, err
		}

		br := bytes.NewReader(ba)
		zr, err := zip.NewReader(br, int64(len(ba)))
		if common.Error(err) {
			return ba, err
		}

		buf := bytes.Buffer{}

		for _, f := range zr.File {
			i, err := f.Open()
			if common.Error(err) {
				return ba, err
			}
			defer func() {
				common.Error(i.Close())
			}()

			_, err = io.Copy(&buf, i)
			if common.Error(err) {
				return buf.Bytes(), err
			}
		}

		return buf.Bytes(), nil
	} else {
		ba, err := hex.DecodeString(this.Content)
		if common.Error(err) {
			return ba, err
		}

		return ba, nil
	}
}

func (this *BinpackFile) UnpackFile(path string) error {
	fileName := filepath.ToSlash(filepath.Join(path, filepath.Base(this.File)))

	common.DebugFunc("Unpack file %s --> %s", this.File, fileName)

	ba, err := this.Unpack()
	if common.Error(err) {
		return err
	}

	err = ioutil.WriteFile(fileName, ba, common.DefaultFileMode)
	if common.Error(err) {
		return err
	}

	return nil
}

func UnpackDir(src string, dest string) error {
	common.DebugFunc("Unpack dir %s --> %s", src, dest)

	for _, file := range BinpackFiles {
		if strings.HasPrefix(file.File, src) {
			fn := filepath.ToSlash(filepath.Clean(filepath.Join(dest, file.File)))

			path := filepath.Dir(fn)
			if !common.FileExists(path) {
				err := os.MkdirAll(path, common.DefaultDirMode)
				if common.Error(err) {
					return err
				}
			}

			err := file.UnpackFile(path)
			if common.Error(err) {
				return err
			}
		}
	}

	return nil
}
