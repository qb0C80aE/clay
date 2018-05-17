if (Data.file_path == undefined) {
  error("form field data[file_path] is undefined");
}

var result = ModelStore.Single("designs/present", "");

if (result[1] != null) {
  error(result[1].Error());
}

result = Conversion.YAMLMarshal(result[0]);

if (result[1] != null) {
  error(result[1].Error());
}

result = Conversion.Bytes(result[0]);

if (result[1] != null) {
  error(result[1].Error());
}

var design = result[0];

result = IO.WriteFile(Data.file_path, design, 0644);

if (result != null) {
  error(result.Error());
}

return "ok";
