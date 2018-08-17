function getEnvironment() {
  var singleResult = ModelStore.Single('environments/1', '');
  if (singleResult[1] != null) {
    error(singleResult[1].Error() + ': ' + Conversion.String(singleResult[0]));
  }
  return singleResult[0];
}

function initGitRepository(environment) {
  var cmdCombinedOutput = Exec.Command(['mkdir', '-p', environment.GitRepositoryURI]).CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  var cmd = Exec.Command(['git', 'status']);
  cmd.Dir = environment.GitRepositoryURI;
  var cmdRun  = cmd.Run();
  if (cmdRun != null) {
    cmd = Exec.Command(['git', 'init']);
    cmd.Dir = environment.GitRepositoryURI;
    cmd.Run();
  }

  cmd = Exec.Command(['git', 'config', '--local', 'user.name', environment.GitUserName]);
  cmd.Dir = environment.GitRepositoryURI;
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  cmd = Exec.Command(['git', 'config', '--local', 'user.email', environment.GitUserEmail]);
  cmd.Dir = environment.GitRepositoryURI;
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }
}

function createTemplateFile(environment) {
  var singleResult = ModelStore.Single('templates/Environment_Template/generation', 'key_parameter=name');
  if (singleResult[1] != null) {
    error(singleResult[1].Error());
  }

  var singleResultBytes = Conversion.Bytes(singleResult[0]);
  if (singleResultBytes[1] != null) {
    error(singleResultBytes[1].Error());
  }

  var writeResult = IO.WriteFile(environment.GitRepositoryURI + '/' + environment.TemplateFileName, singleResultBytes[0], 0644);
  if (writeResult != null) {
    error(result.Error());
  }
}

function createTestRunnerScriptFile(environment) {
  var singleResult = ModelStore.Single('templates/Environment_TestRunnerScript/generation', 'key_parameter=name');
  if (singleResult[1] != null) {
    error(singleResult[1].Error());
  }

  var singleResultBytes = Conversion.Bytes(singleResult[0]);
  if (singleResultBytes[1] != null) {
    error(singleResultBytes[1].Error());
  }

  var writeResult = IO.WriteFile(environment.GitRepositoryURI + '/test_runner.sh', singleResultBytes[0], 0755);
  if (writeResult != null) {
    error(result.Error());
  }
}

function createDesignFile(environment) {
  var singleResult = ModelStore.Single('designs/present', '');
  if (singleResult[1] != null) {
    error(singleResult[1].Error());
  }

  var singleResultJSON = Conversion.JSONMarshal(singleResult[0], '  ');
  if (singleResultJSON[1] != null) {
    error(singleResultJSON[1].Error());
  }

  var singleResultJSONBytes = Conversion.Bytes(singleResultJSON[0]);
  if (singleResultJSONBytes[1] != null) {
    error(singleResultJSONBytes[1].Error());
  }

  var writeResult = IO.WriteFile(environment.GitRepositoryURI + '/' + environment.DesignFileName, singleResultJSONBytes[0], 0644);
  if (writeResult != null) {
    error(result.Error());
  }
}

function createNodeConfigFiles(environment) {
  var cmd = Exec.Command(['rm', '-rf', environment.NodeConfigDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
	var cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

	cmd = Exec.Command(['mkdir', environment.NodeConfigDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  var multiResult = ModelStore.Multi('nodes', 'preloads=node_kind,node_configuration.initialization_template,node_configuration.configuration_template');
  if (multiResult[1] != null) {
    error(multiResult[1].Error());
  }
  var nodes = multiResult[0].Records;

	for (i in nodes) {
    var node = nodes[i];
		var cmd = Exec.Command(['mkdir', '-p', node.Name]);
		cmd.Dir = String.Sprintf('%s/%s', [environment.GitRepositoryURI, environment.NodeConfigDirectoryName]);
		var cmdCombinedOutput = cmd.CombinedOutput();
    if (cmdCombinedOutput[1] != null) {
      error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
    }

    var query = String.Sprintf('p[node_id]=%d', [node.ID]);

		if (!node.Virtual && (Query.Get('init') == 'true')) {
      var templateID = node.NodeKind.InitializationTemplateID;
      if (node.NodeConfiguration.InitializationTemplate.TemplateContent.length > 0) {
        templateID = node.NodeConfiguration.InitializationTemplateID;
      }
      var singleResult = ModelStore.Single('templates/' + templateID + '/generation', query);
      if (singleResult[1] != null) {
        error(singleResult[1].Error());
      }

      singleResultBytes = Conversion.Bytes(singleResult[0]);
      if (singleResultBytes[1] != null) {
        error(singleResultBytes[1].Error());
      }

      var fileName = String.Sprintf('%s/%s/%s/initialize.txt',
        [
          environment.GitRepositoryURI,
          environment.NodeConfigDirectoryName,
          node.Name
        ]
      );

      var writeResult = IO.WriteFile(fileName, singleResultBytes[0], 0644);
      if (writeResult != null) {
        error(result.Error());
      }
		}

    var templateID = node.NodeKind.ConfigurationTemplateID;
    if (node.NodeConfiguration.ConfigurationTemplate.TemplateContent.length > 0) {
      templateID = node.NodeConfiguration.ConfigurationTemplateID;
    }

    var singleResult = ModelStore.Single('templates/' + templateID + '/generation', query);
    if (singleResult[1] != null) {
      error(singleResult[1].Error());
    }

    var singleResultBytes = Conversion.Bytes(singleResult[0]);
    if (singleResultBytes[1] != null) {
      error(singleResultBytes[1].Error());
    }

    var fileName = String.Sprintf('%s/%s/%s/config.txt',
      [
        environment.GitRepositoryURI,
        environment.NodeConfigDirectoryName,
        node.Name
      ]
    );

    var writeResult = IO.WriteFile(fileName, singleResultBytes[0], 0644);
    if (writeResult != null) {
      error(result.Error());
    }
	}
}

function generateTestScript(environment, testScenarioDirectoryName, testScenario) {
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

  var testScenarioParameterConsumedMap = {
    0: {}
  };

  for (i in testScenario.TestScenarioTestStepAssociations) {
    var testScenarioTestStepAssociation = testScenario.TestScenarioTestStepAssociations[i];

    var singleResult = ModelStore.Single('test_steps/' + testScenarioTestStepAssociation.TestStepID, '');
    if (singleResult[1] != null) {
      error(singleResult[1].Error());
    }
    var testStep = singleResult[0];

    singleResult = ModelStore.Single('templates/' + testStep.TemplateID, '');
    if (singleResult[1] != null) {
      error(singleResult[1].Error());
    }
    var template = singleResult[0];

    var multiResult = ModelStore.Multi('template_arguments', 'q[template_id]=' + template.ID);
    if (multiResult[1] != null) {
      error(multiResult[1].Error());
    }
    var templateArguments = multiResult[0].Records;

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

    var singleResult = ModelStore.Single('templates/' + testStep.TemplateID + '/generation', query);
    if (singleResult[1] != null) {
      error(singleResult[1].Error());
    }

    singleResultBytes = Conversion.Bytes(singleResult[0]);
    if (singleResultBytes[1] != null) {
      error(singleResultBytes[1].Error());
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

    var writeResult = IO.WriteFile(fileName, singleResultBytes[0], 0644);
    if (writeResult != null) {
      error(result.Error());
    }
  }

  for (var number in testScenarioParameterMap) {
    var specificNumberTestScenarioParameterMap = testScenarioParameterMap[number];
    var specificNumberTestScenarioParameterConsumedMap = testScenarioParameterConsumedMap[number];
    if (specificNumberTestScenarioParameterConsumedMap == undefined) {
      error('the parameter for the test step ' + number + ' is not used in any test steps');
    }

    for (var key in specificNumberTestScenarioParameterMap) {
      if (specificNumberTestScenarioParameterConsumedMap[key] == undefined) {
        error('the parameter '+ key +' does not exist in step ' + number);
      }
    }
  }
}

function createTestCaseFiles(environment) {
  var cmd = Exec.Command(['rm', '-rf', environment.TestCaseDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
	var cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

	cmd = Exec.Command(['mkdir', environment.TestCaseDirectoryName]);
	cmd.Dir = environment.GitRepositoryURI;
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  var multiResult = ModelStore.Multi('test_scenarios', 'preloads=test_scenario_parameters,test_scenario_test_step_associations');
  if (multiResult[1] != null) {
    error(multiResult[1].Error());
  }

  for (i in multiResult[0].Records) {
    var testScenario = multiResult[0].Records[i];
    var cmd = Exec.Command(['mkdir', environment.TestCaseDirectoryName + '/' + testScenario.Name]);
    cmd.Dir = environment.GitRepositoryURI;
    var cmdCombinedOutput = cmd.CombinedOutput();
    if (cmdCombinedOutput[1] != null) {
      error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
    }

    var testScenarioDirectoryName = environment.GitRepositoryURI + '/' + environment.TestCaseDirectoryName +'/' + testScenario.Name;
    generateTestScript(environment, testScenarioDirectoryName, testScenario);
  }
}

function commit(environment) {
  var cmd = Exec.Command(['git', 'add', '-u']);
	cmd.Dir = environment.GitRepositoryURI;
  var cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  cmd = Exec.Command(['git', 'add', '-A']);
	cmd.Dir = environment.GitRepositoryURI;
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  cmd = Exec.Command(['date']);
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }

  var commitMessage = String.Sprintf('Automatic commit at %s', Conversion.String(cmdCombinedOutput[0]));
  cmd = Exec.Command(['git', 'commit', '-m', commitMessage]);
	cmd.Dir = environment.GitRepositoryURI;
  cmdCombinedOutput = cmd.CombinedOutput();
  if (cmdCombinedOutput[1] != null) {
    error(cmdCombinedOutput[1].Error() + ': ' + Conversion.String(cmdCombinedOutput[0]));
  }
}

function main() {
  var environment = getEnvironment();
  initGitRepository(environment);
  createTemplateFile(environment);
  createTestRunnerScriptFile(environment);
  createDesignFile(environment);
  createNodeConfigFiles(environment);
  createTestCaseFiles(environment);
  commit(environment);
}

main();
return 'ok';
