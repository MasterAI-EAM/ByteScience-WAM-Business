package conf

var StepNameData = map[string]string{
	"C1":  "树脂 25℃粘度/mPa·s",
	"C2":  "环氧当量/g/eq",
	"C3":  "固化剂 25℃粘度/mPa·s",
	"C4":  "胺值/mg KOH/g",
	"D1":  "25℃混合粘度/mPa·s",
	"D2":  "25℃可使用时间/min",
	"D3":  "30℃混合粘度/mPa·s",
	"D4":  "30℃可使用时间/min",
	"D5":  "35℃混合粘度/mPa·s",
	"D6":  "35℃可使用时间/min",
	"E1":  "拉伸强度Mpa",
	"E2":  "拉伸模量Mpa",
	"E3":  "断裂延伸率%",
	"E4":  "弯曲强度Mpa",
	"E5":  "弯曲模量Mpa",
	"E6":  "压缩强度Mpa",
	"E7":  "冲击韧性KJ/m2",
	"E8":  "Tg/'C",
	"I1":  "0°拉伸强度Mpa",
	"I2":  "0°拉伸模量Gpa",
	"I3":  "0°拉伸应变%",
	"I4":  "泊松比",
	"I5":  "90°拉伸强度Mpa",
	"I6":  "90°拉伸模量Gpa",
	"I7":  "90°拉伸应变%",
	"I8":  "0°压缩强度Mpa",
	"I9":  "0°压缩模量Gpa",
	"I10": "0°压缩应变%",
	"I11": "90°压缩强度Mpa",
	"I12": "90°模量Gpa",
	"I13": "90°压缩应变%",
	"I14": "V剪强度Mpa",
	"I15": "V剪模量Gpa",
	"I16": "V剪应变%",
}

const (
	Login            = "login"
	ChangPassword    = "changPassword"
	ImportExperiment = "importExperiment"
	AddExperiment    = "addExperiment"
	EditExperiment   = "editExperiment"
	DeleteExperiment = "deleteExperiment"
	AddRecipe        = "addRecipe"
	EditRecipe       = "editRecipe"
	DeleteRecipe     = "deleteRecipe"
)
