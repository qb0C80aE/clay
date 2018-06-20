function getEnvironment() {
  var result = ModelStore.Single('environments/1', '');
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }
  return result[0];
}

function initGitRepository(environment) {
  var result = Exec.Command(['mkdir', '-p', environment.GitRepositoryURI]).CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  var cmd = Exec.Command(['git', 'status']);
  cmd.Dir = environment.GitRepositoryURI;
  result = cmd.Run();
  if (result != null) {
    cmd = Exec.Command(['git', 'init']);
    cmd.Dir = environment.GitRepositoryURI;
    cmd.Run();
  }

  cmd = Exec.Command(['git', 'config', '--local', 'user.name', environment.GitUserName]);
  cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  cmd = Exec.Command(['git', 'config', '--local', 'user.email', environment.GitUserEmail]);
  cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }
}

function unzip(environment) {
  if (Data.zip_file == undefined) {
    error('form field data[zip_file] is undefined');
  }

  if (Data.notification_url == undefined) {
    error('form field data[notification_url] is undefined');
  }

  var cmd = Exec.Command(['rm', '-rf', 'uploaded']);
	cmd.Dir = environment.GitRepositoryURI;
	result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  cmd = Exec.Command(['mkdir', '-p', 'uploaded']);
  cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  var fileName = String.Sprintf('%s/%s/uploaded.zip', environment.GitRepositoryURI, 'uploaded');
  result = IO.WriteFile(fileName, Data.zip_file, 0644);
  if (result != null) {
    error(result.Error());
  }

  cmd = Exec.Command(['unzip', 'uploaded.zip']);
  cmd.Dir = String.Sprintf('%s/%s', environment.GitRepositoryURI, 'uploaded');
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  fileName = String.Sprintf('%s/%s/notification_url.txt', environment.GitRepositoryURI, 'uploaded');
  result = Conversion.Bytes(Data.notification_url);
  if (result[1] != null) {
    error(result[1].Error());
  }
  result = IO.WriteFile(fileName, result[0], 0644);
  if (result != null) {
    error(result.Error());
  }
}

function updateDesign(environment) {
  cmd = Exec.Command(['curl', '-X', 'PUT', '-H', 'Content-Type: application/json', 'http://localhost:8080/designs/present', '-d', '@design.json']);
  cmd.Dir = String.Sprintf('%s/%s', environment.GitRepositoryURI, 'uploaded');
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }
}

function createTemplateFile(environment) {
  var result = ModelStore.Single('templates/environment_template/generation', 'key_parameter=name');
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = Conversion.Bytes(result[0]);
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = IO.WriteFile(environment.GitRepositoryURI + '/' + environment.TemplateFileName, result[0], 0644);
  if (result != null) {
    error(result.Error());
  }
}

function createTestRunnerScriptFile(environment) {
  var result = ModelStore.Single('templates/environment_test_runner_script/generation', 'key_parameter=name');
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = Conversion.Bytes(result[0]);
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = IO.WriteFile(environment.GitRepositoryURI + '/test_runner.sh', result[0], 0644);
  if (result != null) {
    error(result.Error());
  }
}

function createDesignFile(environment) {
  var result = ModelStore.Single('designs/present', '');
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = Conversion.JSONMarshal(result[0], '  ');
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = Conversion.Bytes(result[0]);
  if (result[1] != null) {
    error(result[1].Error());
  }

  result = IO.WriteFile(environment.GitRepositoryURI + '/' + environment.DesignFileName, result[0], 0644);
  if (result != null) {
    error(result.Error());
  }
}

function createNodeConfigFiles(environment) {
  var cmd = Exec.Command(['rm', '-rf', environment.NodeConfigDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
	result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

	cmd = Exec.Command(['mkdir', environment.NodeConfigDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  var result = ModelStore.Multi('nodes', 'preloads=node_kind,node_configuration.initialization_template,node_configuration.configuration_template');
  if (result[1] != null) {
    error(result[1].Error());
  }
  var nodes = result[0].Records;

	for (i in nodes) {
    var node = nodes[i];
		var cmd = Exec.Command(['mkdir', '-p', node.Name]);
		cmd.Dir = String.Sprintf('%s/%s', [environment.GitRepositoryURI, environment.NodeConfigDirectoryName]);
		result = cmd.CombinedOutput();
    if (result[1] != null) {
      error(result[1].Error() + ': ' + Conversion.String(result[0]));
    }

    var query = String.Sprintf('node_id=%d', [node.ID]);

		if (!node.Virtual && (Query.Get('init') == 'true')) {
      var templateID = node.NodeKind.InitializationTemplateID;
      if (node.NodeConfiguration.InitializationTemplate.TemplateContent.length > 0) {
        templateID = node.NodeConfiguration.InitializationTemplateID;
      }
      var result = ModelStore.Single('templates/' + templateID + '/generation', query);
      if (result[1] != null) {
        error(result[1].Error());
      }

      result = Conversion.Bytes(result[0]);
      if (result[1] != null) {
        error(result[1].Error());
      }

      var fileName = String.Sprintf('%s/%s/%s/initialize.txt',
        [
          environment.GitRepositoryURI,
          environment.NodeConfigDirectoryName,
          node.Name
        ]
      );

      result = IO.WriteFile(fileName, result[0], 0644);
      if (result != null) {
        error(result.Error());
      }
		}

    var result = ModelStore.Single('templates/' + node.NodeKind.ConfigurationTemplateID + '/generation', query);
    if (result[1] != null) {
      error(result[1].Error());
    }

    result = Conversion.Bytes(result[0]);
    if (result[1] != null) {
      error(result[1].Error());
    }

    var fileName = String.Sprintf('%s/%s/%s/config.txt',
      [
        environment.GitRepositoryURI,
        environment.NodeConfigDirectoryName,
        node.Name
      ]
    );

    result = IO.WriteFile(fileName, result[0], 0644);
    if (result != null) {
      error(result.Error());
    }
	}
}

function generateTestScript(testScenarioDirectoryName, testScenario) {
  var testScenarioParameterMap = {
    0: {}
  };

  for (i in testScenario.TestScenarioParameters) {
    var testScenarioParameter = testScenario.TestScenarioParameters[i];
    if (testScenarioParameterMap[testScenarioParameter.TestStepNumber] == undefined) {
      testScenarioParameterMap[testScenarioParameter.TestStepNumber] = {};
    }
    testScenarioParameterMap[testScenarioParameter.TestStepNumber][testScenarioParameter.Name] = testScenarioParameter.Value;
  }

  testScenarioParameterConsumedMap = {
    0: {}
  };

  for (i in testScenario.TestScenarioTestStepAssociations) {
    var testScenarioTestStepAssociation = testScenario.TestScenarioTestStepAssociations[i];

    var result = ModelStore.Single('test_steps/' + testScenarioTestStepAssociation.TestStepID, '');
    if (result[1] != null) {
      error(result[1].Error());
    }
    var testStep = result[0];

    var result = ModelStore.Single('templates/' + testStep.TemplateID, '');
    if (result[1] != null) {
      error(result[1].Error());
    }
    var template = result[0];

    var result = ModelStore.Multi('template_arguments', 'q[template_id]=' + template.ID);
    if (result[1] != null) {
      error(result[1].Error());
    }
    var templateArguments = result[0].Records;

    var globalParameterMap = testScenarioParameterMap[0];
    var localParameterMap = globalParameterMap;
    if (testScenarioParameterMap[testScenarioTestStepAssociation.Number] != undefined) {
      localParameterMap = testScenarioParameterMap[testScenarioTestStepAssociation.Number];
    }

    var templateArgumentParameterMap = {};
    for (j in templateArguments) {
      var templateArgument = templateArguments[j];
      var templateParameterValue = globalParameterMap[templateArgument.Name];
      if (templateParameterValue != undefined) {
        templateArgumentParameterMap[templateArgument.Name] = templateParameterValue;
        testScenarioParameterConsumedMap[0][templateArgument.Name] = true;
      }

      templateParameterValue = localParameterMap[templateArgument.Name];
      if (templateParameterValue != undefined) {
        templateArgumentParameterMap[templateArgument.Name] = templateParameterValue;

        if (testScenarioParameterConsumedMap[testScenarioTestStepAssociation.Number] == undefined) {
          testScenarioParameterConsumedMap[testScenarioTestStepAssociation.Number] = {};
        }

        testScenarioParameterConsumedMap[testScenarioTestStepAssociation.Number][templateArgument.Name] = true;
      }
    }

    var query = '';
    for (key in templateArgumentParameterMap) {
      if (query.length > 0) {
        query = query + '&';
      }
      query = query + 'p[' + key + ']=' + templateArgumentParameterMap[key];
    }

    var result = ModelStore.Single('templates/' + testStep.TemplateID + '/generation', query);
    if (result[1] != null) {
      error(result[1].Error());
    }

    result = Conversion.Bytes(result[0]);
    if (result[1] != null) {
      error(result[1].Error());
    }

    var fileName = String.Sprintf('%s/%s/%s/%s_%s',
      [
        environment.GitRepositoryURI,
        environment.TestCaseDirectoryName,
        testScenario.Name,
        String.Sprintf('%07d', testScenarioTestStepAssociation.Number),
        testStep.Name
      ]
    );
    result = IO.WriteFile(fileName, result[0], 0644);
    if (result != null) {
      error(result.Error());
    }
  }

  for (number in testScenarioParameterMap) {
    var specificNumberTestScenarioParameterMap = testScenarioParameterMap[number];
    var specificNumberTestScenarioParameterConsumedMap = testScenarioParameterConsumedMap[number];
    if (specificNumberTestScenarioParameterConsumedMap == undefined) {
      error('the parameter for the test step ' + number + ' is not used in any test steps');
    }

    for (key in specificNumberTestScenarioParameterMap) {
      if (specificNumberTestScenarioParameterConsumedMap[key] == undefined) {
        error('the parameter '+ key +' does not exist in step ' + number);
      }
    }
  }
}

function createTestCaseFiles(environment) {
  var cmd = Exec.Command(['rm', '-rf', environment.TestCaseDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
	result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

	cmd = Exec.Command(['mkdir', environment.TestCaseDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  result = ModelStore.Multi('test_scenarios', 'preloads=test_scenario_parameters,test_scenario_test_step_associations');
  if (result[1] != null) {
    error(result[1].Error());
  }

  for (i in result[0].Records) {
    var testScenario = result[0].Records[i];
    cmd = Exec.Command(['mkdir', environment.TestCaseDirectoryName + '/' + testScenario.Name]);
    cmd.Dir = environment.GitRepositoryURI;
    result = cmd.CombinedOutput();
    if (result[1] != null) {
      error(result[1].Error() + ': ' + Conversion.String(result[0]));
    }

    var testScenarioDirectoryName = environment.GitRepositoryURI + '/' + environment.TestCaseDirectoryName +'/' + testScenario.Name;
    generateTestScript(testScenarioDirectoryName, testScenario);
  }
}

function commit(environment) {
  cmd = Exec.Command(['git', 'add', '-u']);
	cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  cmd = Exec.Command(['git', 'add', '-A']);
	cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  cmd = Exec.Command(['date']);
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }

  var commitMessage = String.Sprintf('Automatic commit at %s', Conversion.String(result[0]));
  cmd = Exec.Command(['git', 'commit', '-m', commitMessage]);
	cmd.Dir = environment.GitRepositoryURI;
  result = cmd.CombinedOutput();
  if (result[1] != null) {
    error(result[1].Error() + ': ' + Conversion.String(result[0]));
  }
}

function generate(environment) {
  createTemplateFile(environment);
  createTestRunnerScriptFile(environment);
  createDesignFile(environment);
  createNodeConfigFiles(environment);
  createTestCaseFiles(environment);
}

function main() {
  var environment = getEnvironment();
  initGitRepository(environment);
  unzip(environment);
  updateDesign(environment);
  generate(environment);
  commit(environment);
}

main();
return 'ok';
