//  Copyright DataWiseHQ/grule-rule-engine Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package examples

import (
	"github.com/DataWiseHQ/grule-rule-engine/ast"
	"github.com/DataWiseHQ/grule-rule-engine/builder"
	"github.com/DataWiseHQ/grule-rule-engine/engine"
	"github.com/DataWiseHQ/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Fact struct {
	NetAmount float32
	Distance  int32
	Duration  int32
	Result    bool
}

const duplicateRulesWithDiffSalience = `
rule  DuplicateRule1  "Duplicate Rule 1"  salience 5 {
when
(Fact.Distance > 5000  &&   Fact.Duration > 120) && (Fact.Result == false)
Then
   Fact.NetAmount=143.320007;
   Fact.Result=true;
}

rule  DuplicateRule2  "Duplicate Rule 2"  salience 6 {
when
(Fact.Distance > 5000  &&   Fact.Duration > 120) && (Fact.Result == false)
Then
   Fact.NetAmount=143.320007;
   Fact.Result=true;
}


rule  DuplicateRule3  "Duplicate Rule 3"  salience 7 {
when
(Fact.Distance > 5000  &&   Fact.Duration > 120) && (Fact.Result == false)
Then
   Fact.NetAmount=143.320007;
   Fact.Result=true;
}


rule  DuplicateRule4  "Duplicate Rule 4"  salience 8 {
when
(Fact.Distance > 5000  &&   Fact.Duration > 120) && (Fact.Result == false)
Then
   Fact.NetAmount=143.320007;
   Fact.Result=true;
}


rule  DuplicateRule5  "Duplicate Rule 5"  salience 9 {
when
(Fact.Distance > 5000  &&   Fact.Duration == 120) && (Fact.Result == false)
Then
   Output.NetAmount=143.320007;
   Fact.Result=true;
}`

func TestGruleEngine_FetchMatchingRules_Having_Diff_Salience(t *testing.T) {
	//Given
	fact := &Fact{
		Distance: 6000,
		Duration: 121,
	}
	dctx := ast.NewDataContext()
	err := dctx.Add("Fact", fact)
	assert.NoError(t, err)
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err = rb.BuildRuleFromResource("conflict_rules_test", "0.1.1", pkg.NewBytesResource([]byte(duplicateRulesWithDiffSalience)))
	assert.NoError(t, err)
	kb, err := lib.NewKnowledgeBaseInstance("conflict_rules_test", "0.1.1")
	assert.NoError(t, err)

	//When
	e := engine.NewGruleEngine()
	ruleEntries, err := e.FetchMatchingRules(dctx, kb)
	assert.NoError(t, err)

	//Then
	assert.Equal(t, 4, len(ruleEntries))
	assert.Equal(t, 8, ruleEntries[0].Salience)
	assert.Equal(t, 7, ruleEntries[1].Salience)
	assert.Equal(t, 6, ruleEntries[2].Salience)
	assert.Equal(t, 5, ruleEntries[3].Salience)
}
