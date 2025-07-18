package wwise

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadTranslationCSV(versionTable string, nameTable string) error {
	versionFile, err := os.Open(versionTable)
	if err != nil {
		return fmt.Errorf("Failed to open %s version table file: %w", versionTable, err)
	}
	nameFile, err := os.Open(nameTable)
	if err != nil {
		return fmt.Errorf("Failed to open %s name table file: %w", nameTable, err)
	}
	versionReader := csv.NewReader(versionFile)
	nameReader := csv.NewReader(nameFile)

	versionLUTs := make(map[int]map[PropType]uint8)
	versions := []int{}

	nameLUT := make(map[PropType]string)
	nameGuardLUT := make(map[string]PropType)
	enumNameLUT := make(map[PropType]string)
	enumNameGuardLUT := make(map[string]PropType)
	evArr := []PropType{}

	j := PropType(0)
	for {
		rows, err := nameReader.Read()
		if err != nil {
			if err == io.EOF {
				if j == 0 {
					return fmt.Errorf("Name table %s provide no enum", nameTable)
				}
				break
			}
			return fmt.Errorf("Failed to read a row from name table %s: %w", nameTable, err)
		}

		rows[0] = "T" + rows[0]
		if _, in := enumNameLUT[j]; in {
			return fmt.Errorf("Enum value %d already exists in the table", j)
		}
		enumNameLUT[j] = rows[0]
		if _, in := enumNameGuardLUT[rows[0]]; in {
			return fmt.Errorf("Enum name %s already exists in the table", rows[0])
		}
		enumNameGuardLUT[rows[0]] = j

		if _, in := nameLUT[j]; in {
			return fmt.Errorf("Enum value %d already exists in the table", j)
		}
		nameLUT[j] = rows[1]
		if _, in := nameGuardLUT[rows[1]]; in {
			return fmt.Errorf("Enum map name %s already exists in the table", rows[1])
		}
		nameGuardLUT[rows[1]] = j
		evArr = append(evArr, j)
		j += 1
	}
	if !(len(nameGuardLUT) == len(nameLUT) && len(enumNameLUT) == len(enumNameGuardLUT) && len(nameLUT) == len(enumNameLUT) && len(evArr) == len(enumNameLUT) && len(nameLUT) == int(j)) {
		return fmt.Errorf("# of key-value pair in enum name LUT and name LUT is mistach.")
	}

	versionRow, err := versionReader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain version values from version table %s: %w", versionTable, err)
	}
	if len(versionRow) <= 1 {
		return fmt.Errorf("First row of version table %s only has one column", versionTable)
	}
	for i := range versionRow {
		if i == 0 {
			continue
		}
		version, err := strconv.Atoi(versionRow[i])
		if err != nil {
			return fmt.Errorf("Failed to parse version value at column %d", i)
		}
		if _, in := versionLUTs[version]; in {
			return fmt.Errorf("Version %d already exists in the table", version)
		}
		versionLUTs[version] = make(map[PropType]uint8, j)
		versions = append(versions, version)
	}
	k := PropType(0)
	for {
		rows, err := versionReader.Read()
		if err != nil {
			if err == io.EOF {
				if k == 0 {
					return fmt.Errorf("Version table %s provide no version translations", versionTable)
				}
				break
			}
		}
		if len(rows) - 1 != len(versionLUTs) {
			return fmt.Errorf("# of columns does not equal to # of version at row %d", k + 1)
		}

		rows[0] = "T" + rows[0]
		ev, in := enumNameGuardLUT[rows[0]]; 
		if !in {
			return fmt.Errorf("Enum %s does not exist in name table %s", rows[0], nameTable)
		}
		if ev != k {
			return fmt.Errorf("Enum value for %s mismatch from version table %s", rows[0], versionTable)
		}

		for i := range rows {
			if i == 0 {
				continue
			}
			version := versions[i - 1]
			versionLUT, in := versionLUTs[version]
			if !in {
				return fmt.Errorf("Version %d does not have translation table from version table", version)
			}

			if _, in := versionLUT[ev]; in {
				return fmt.Errorf("Enum value %d already has a translation", ev)
			}
			versionPropID, err := strconv.ParseUint(rows[i], 16, 8)
			if err != nil {
				return fmt.Errorf("Failed to parse version specific property ID: %w", err)
			}
			versionLUT[ev] = uint8(versionPropID)
		}

		k += 1
	}

	if k != j {
		return fmt.Errorf("# of enum from version table doesn't equal # of enum from name table")
	}

	var pyBuilder strings.Builder
	var builder strings.Builder
	builder.WriteString("package wwise\n\n")
	for version, versionLUT := range versionLUTs {
		pyBuilder.WriteString(fmt.Sprintf("ForwardTranslationV%d = {\n", version))
		builder.WriteString(fmt.Sprintf("var ForwardTranslationV%d = map[PropType]uint8{\n", version))
		for _, ev := range evArr {
			enumName, in := enumNameLUT[ev]
			if !in {
				return fmt.Errorf("Enum value %d does not have a enum name", ev)
			}
			versionPropId, in := versionLUT[ev]
			if !in {
				return fmt.Errorf("Enum value %d does not have property ID for version %d", ev, version)
			}
			if versionPropId == 0xFF {
				continue
			}
			pyBuilder.WriteString(fmt.Sprintf("    %s: %d,\n", enumName, versionPropId))
			builder.WriteString(fmt.Sprintf("    %s: %d,\n", enumName, versionPropId))
		}
		builder.WriteString("}\n")
		pyBuilder.WriteString("}\n")

		pyBuilder.WriteString(fmt.Sprintf("InverseTranslationV%d = {\n", version))
		builder.WriteString(fmt.Sprintf("var InverseTranslationV%d = map[uint8]PropType{\n", version))
		for _, ev := range evArr {
			enumName, in := enumNameLUT[ev]
			if !in {
				return fmt.Errorf("Enum value %d does not have a enum name", ev)
			}
			versionPropId, in := versionLUT[ev]
			if !in {
				return fmt.Errorf("Enum value %d does not have property ID for version %d", ev, version)
			}
			if versionPropId == 0xFF {
				continue
			}
			pyBuilder.WriteString(fmt.Sprintf("    %d: %s,\n", versionPropId, enumName))
			builder.WriteString(fmt.Sprintf("    %d: %s,\n", versionPropId, enumName))
		}
		builder.WriteString("}\n\n")
		pyBuilder.WriteString("}\n\n")
	}

	var pyNameBuilder strings.Builder
	var nameBuilder strings.Builder

	builder.WriteString("const (\n")
	pyNameBuilder.WriteString("TranslateName = {\n")
	nameBuilder.WriteString("var TranslateName = map[PropType]string{\n")
	for _, ev := range evArr {
		enumName, in := enumNameLUT[ev]
		if !in {
			return fmt.Errorf("Enum value %d does not have a enum name", ev)
		}
		name, in := nameLUT[ev]
		if !in {
			return fmt.Errorf("Enum value %d does not have a enum string name", ev)
		}
		builder.WriteString(fmt.Sprintf("    %s PropType = %d\n", enumName, ev))
		pyBuilder.WriteString(fmt.Sprintf("%s = %d\n", enumName, ev))
		nameBuilder.WriteString(fmt.Sprintf("    %s: \"%s\",\n", enumName, name))
		pyNameBuilder.WriteString(fmt.Sprintf("    %s: \"%s\",\n", enumName, name))
	}

	builder.WriteString(")\n\n")
	nameBuilder.WriteString("}\n")
	pyBuilder.WriteString("\n\n")
	pyNameBuilder.WriteString("}\n")

	out := builder.String() + nameBuilder.String()

	if err := os.WriteFile("prop_translation_lut.go", []byte(out), 0666); err != nil {
		return err
	}

	out = pyBuilder.String() + pyNameBuilder.String()
	return os.WriteFile("prop_translation_lut.py", []byte(out), 0666)
}
