package main

import (
	"log"
	"os"
	"text/template"
)

// Node representa un equipo de red en Containerlab
type Node struct {
	Name     string
	Kind     string
	Image    string
	MgmtIPv4 string
}

// Link representa una conexión punto a punto o multipunto
type Link struct {
	Endpoints []string
}

// Topology representa la estructura completa del laboratorio
type Topology struct {
	Name       string
	MgmtSubnet string
	Nodes      []Node
	Links      []Link
}

func main() {
	// Definición de la topología modernizada del IXP
	topo := Topology{
		Name:       "ixp-lab",
		MgmtSubnet: "10.254.0.0/24",
		Nodes: []Node{
			{Name: "r1", Kind: "vr-vmx", Image: "vrnetlab/vr-vmx:latest", MgmtIPv4: "10.254.0.11"},
			{Name: "r3", Kind: "vr-vmx", Image: "vrnetlab/vr-vmx:latest", MgmtIPv4: "10.254.0.13"},
			{Name: "r4", Kind: "vr-vmx", Image: "vrnetlab/vr-vmx:latest", MgmtIPv4: "10.254.0.14"},
			// Nodos GoBGP modernizados con Alpine y FRR
			{Name: "g2", Kind: "linux", Image: "alpine-gobgp-frr:latest", MgmtIPv4: "10.254.0.102"},
			{Name: "rs", Kind: "linux", Image: "alpine-gobgp:latest", MgmtIPv4: "10.254.0.150"},
			{Name: "util", Kind: "linux", Image: "alpine:3.18", MgmtIPv4: "10.254.0.250"},
			// Switch bridge para simular el Punto de Intercambio (L2)
			{Name: "ix-switch", Kind: "bridge", Image: "", MgmtIPv4: ""},
		},
		Links: []Link{
			// r1 - r4 (eBGP directo)
			{Endpoints: []string{"r1:eth1", "r4:eth1"}},
			// r1 - g2 (iBGP & OSPF)
			{Endpoints: []string{"r1:eth2", "g2:eth1"}},
			// IX Switch (Route Server + Peers)
			{Endpoints: []string{"g2:eth2", "ix-switch:eth1"}},
			{Endpoints: []string{"rs:eth1", "ix-switch:eth2"}},
			{Endpoints: []string{"r3:eth1", "ix-switch:eth3"}},
			{Endpoints: []string{"r1:eth3", "ix-switch:eth4"}},
		},
	}

	// 1. Cargar la plantilla base
	tmpl, err := template.ParseFiles("topology.tmpl")
	if err != nil {
		log.Fatalf("Error leyendo plantilla topology.tmpl: %v\n", err)
	}

	// 2. Crear o sobrescribir el archivo destino YAML de Containerlab
	f, err := os.Create("ixp-lab.clab.yml")
	if err != nil {
		log.Fatalf("Error creando el archivo clab.yml: %v\n", err)
	}
	defer f.Close()

	// 3. Renderizar y aplicar los datos a la plantilla
	log.Println("Generando topología ixp-lab.clab.yml con los datos en Go...")
	err = tmpl.Execute(f, topo)
	if err != nil {
		log.Fatalf("Error renderizando la plantilla: %v\n", err)
	}

	log.Println("¡Archivo generado existosamente! Listo para ejecutar: sudo clab deploy -t ixp-lab.clab.yml")
}
