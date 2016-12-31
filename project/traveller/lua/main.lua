WorldLoop = function()
end 

NodeRadar                   = Node("CanvasRadar")
NodeTerminalMain            = Node("TerminalMain")
NodeParInfo                 = Node("ParInfo")
NodeParNewestMsg            = Node("ParNewestMsg")
NodeParGUserSpaceshipStatus = Node("ParGUserSpaceshipStatus")
NodeParGUserSpaceshipCabin  = Node("ParGUserSpaceshipCabin")
NodeGaugeFuel               = Node("GaugeFuel")
NodeGaugeLife               = Node("GaugeLife")
NodeTabpaneMain             = Node("TabpaneMain")

NodeRadar:SetAttribute("height", tostring(WindowHeight()-NodeTerminalMain:Height()))

GRadar         = NewRadar()
GTerminal      = NewTerminal()
GWorld         = NewWorld()
GUserSpaceship = GetSpaceshipFromDB(1)
GUserSpaceship:RefreshNodeParGUserSpaceshipStatus()
GUserSpaceship:RefreshGaugeLife()
GUserSpaceship:RefreshGaugeFuel()

NodeParGUserSpaceshipStatus:SetAttribute("borderlabel", string.format(" %s状态 ", GUserSpaceship.Info.Name))

NodeTerminalMain:SetActive()
GUserSpaceship:UpdateFuel(0)
GUserSpaceship:UpdateLife(0)
GWorld:LoopEvent()

--[[
SetTimeout(200, function()
    GTerminal.CmdExcuter["/planet"]:ShowPlanetDetail()
end)
]]
