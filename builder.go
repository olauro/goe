package goe

import (
	"fmt"
	"strings"
)

const (
	querySELECT int8 = 1
	queryINSERT int8 = 2
	queryUPDATE int8 = 3
)

var (
	statementSELECT = &statement{
		keyword: "SELECT",
		tip:     writeDML,
	}
	statementFROM = &statement{
		keyword: "FROM",
		tip:     writeDML,
	}
	statementWHERE = &statement{
		keyword: "WHERE",
		tip:     writeMIDDLE,
	}
)

type builder struct {
	sql       *strings.Builder
	args      []string
	brs       []*booleanResult
	queue     *statementQueue
	tables    *statementQueue
	pks       *pkQueue
	queryType int8
}

func createBuilder(qt int8) *builder {
	return &builder{
		sql:       &strings.Builder{},
		queue:     createStatementQueue(),
		tables:    createStatementQueue(),
		queryType: qt,
		pks:       createPkQueue()}
}

type statement struct {
	keyword string
	tip     int8
}

func createStatement(k string, t int8) *statement {
	return &statement{keyword: k, tip: t}
}

func (b *builder) buildSelect(addrMap map[string]any) {
	//TODO: Set a drive type to share stm
	b.queue.add(statementSELECT)

	//TODO Better Query
	for _, v := range b.args {
		switch atr := addrMap[v].(type) {
		case *att:
			b.queue.add(createStatement(atr.name, writeATT))
			b.tables.add(createStatement(atr.pk.table, writeTABLE))

			//TODO: Add a list pk?
			b.pks.add(atr.pk)
		case *pk:
			b.queue.add(createStatement(atr.name, writeATT))
			b.tables.add(createStatement(atr.table, writeTABLE))

			//TODO: Add a list pk?
			b.pks.add(atr)
		}
	}

	b.queue.add(statementFROM)
}

func (b *builder) buildSql() {
	switch b.queryType {
	case querySELECT:
		b.buildTables()
		b.buildWhere()
		writeSelect(b.sql, b.queue)
	case queryINSERT:
		break
	case queryUPDATE:
		break
	}
}

func (b *builder) buildWhere() {
	if len(b.brs) == 0 {
		return
	}
	b.queue.add(statementWHERE)
	for _, br := range b.brs {
		switch br.tip {
		case EQUALS:
			b.queue.add(createStatement(fmt.Sprintf("%v = %v", br.arg, br.value), 0))
		}
	}
}

func (b *builder) buildTables() {
	b.queue.add(b.tables.get())
	if b.tables.size >= 1 {
		for table := b.tables.get(); table != nil; {
			buildJoins(table, b.pks, b.queue)
			table = b.tables.get()
		}
	}
}

func buildJoins(table *statement, pks *pkQueue, stQueue *statementQueue) {
	for pk := pks.get(); pk != nil; {
		if pk.table != table.keyword {
			switch fk := pk.fks[table.keyword].(type) {
			case *manyToOne:
				if fk.hasMany {
					stQueue.add(
						createStatement(fmt.Sprintf("inner join %v on (%v = %v)", table.keyword, pk.name, fk.id), writeJOIN),
					)
				} else {
					stQueue.add(
						createStatement(fmt.Sprintf("inner join %v on (%v = %v)", table.keyword, pks.findPk(table.keyword).name, fk.id), writeJOIN),
					)
				}
			case *manyToMany:
				if !pk.skipFlag {
					stQueue.add(
						createStatement(fmt.Sprintf("inner join %v on (%v = %v)", fk.table, pk.name, fk.ids[pk.table]), writeJOIN),
					)
					stQueue.add(
						createStatement(
							fmt.Sprintf(
								"inner join %v on (%v = %v)",
								table.keyword, fk.ids[table.keyword],
								pks.findPk(table.keyword).name), writeJOIN,
						),
					)
				}
			}
		}
		pk = pks.get()
	}
}
