local _Radar = {}
local _mtRadar = {__index = _Radar} 

function NewRadar()
    local Radar = setmetatable({}, _mtRadar)
    Radar.ScreenPlanets = {}
    Radar.CursorScreenPosition = {X=GetIntPart(NodeRadar:Width()/2), Y=GetIntPart(NodeRadar:Height()/2)}
    Radar.KeyPressStrForMove = ""
    Radar.FocusPlanet= nil

    Radar.KeyPressSig = NodeRadar:RegisterKeyPressHandler(function(nodePointer, keyStr)
        Radar:KeyPressHandle(nodePointer, keyStr)
    end)
    Radar.ActiveModeSig = NodeRadar:RegisterLuaActiveModeHandler(function(nodePointer)
        Radar:ActiveMode(nodePointer)
    end)
    return Radar
end

function _Radar.ActiveMode(self, nodePointer)
    self:renewCursor()
end

function _Radar.RefreshParInfo(self)
    if nil ~= self.FocusTarget then
        self.FocusTarget.ColorBg = ""
    end

    self.FocusTarget = nil

    if self.CursorScreenPosition.X == GUserSpaceship.ScreenPosition.X and
        self.CursorScreenPosition.Y == GUserSpaceship.ScreenPosition.Y then
        self.FocusTarget = GUserSpaceship
        self.FocusTarget.ColorBg = "white"
        NodeInputTextNamePlanet:SetText(GUserSpaceship.Info.Name)
        NodeInputTextNamePlanet:SetAttribute("borderlabel", "飞船")
        NodeParInfo:SetText(string.format([[
X: %d
Y: %d]], GUserSpaceship.Info.Position.X, GUserSpaceship.Info.Position.Y))
        return
    end

    local planet = self.ScreenPlanets[PointToStr(self.CursorScreenPosition)]
    if nil ~= planet then
        self.FocusTarget = planet
        self.FocusTarget.ColorBg = "white"
        NodeInputTextNamePlanet:SetText(planet.Info.Name)
        NodeInputTextNamePlanet:SetAttribute("borderlabel", "星球")
        NodeParInfo:SetText(string.format([[
X: %d
Y: %d
资源: %d]], planet.Info.Position.X, planet.Info.Position.Y, planet.Info.Resource))
        return
    end

    if nil == self.FocusTarget then
        NodeInputTextNamePlanet:SetText("")
        NodeInputTextNamePlanet:SetAttribute("borderlabel", "")
        NodeParInfo:SetText("")
    end
end

function _Radar.KeyPressHandle(self, nodePointer, keyStr)
    if "<enter>" == keyStr then
        if "" == GTerminal.CurrentCommand then
        end
    end

    local movePosition = "no"
    local isMove = false
    if "<left>" == keyStr then
        isMove = true
        self.CursorScreenPosition.X = self.CursorScreenPosition.X - 1
    elseif "<right>" == keyStr then
        isMove = true
        self.CursorScreenPosition.X = self.CursorScreenPosition.X + 1
    elseif "<up>" == keyStr then
        isMove = true
        self.CursorScreenPosition.Y = self.CursorScreenPosition.Y - 1
    elseif "<down>" == keyStr then
        isMove = true
        self.CursorScreenPosition.Y = self.CursorScreenPosition.Y + 1
    end

    self.KeyPressStrForMove = self.KeyPressStrForMove .. keyStr 
    local step
    if "h" == keyStr then
        isMove = true
        step = tonumber(string.sub(self.KeyPressStrForMove, 0, -2))
        if "number" == type(step) then
        else 
            step = 1
        end
        self.CursorScreenPosition.X = self.CursorScreenPosition.X - step
    elseif "l" == keyStr then
        isMove = true
        step = tonumber(string.sub(self.KeyPressStrForMove, 0, -2))
        if "number" == type(step) then
        else 
            step = 1
        end
        self.CursorScreenPosition.X = self.CursorScreenPosition.X + step
    elseif "k" == keyStr then
        isMove = true
        step = tonumber(string.sub(self.KeyPressStrForMove, 0, -2))
        if "number" == type(step) then
        else 
            step = 1
        end
        self.CursorScreenPosition.Y = self.CursorScreenPosition.Y - step
    elseif "j" == keyStr then
        isMove = true
        step = tonumber(string.sub(self.KeyPressStrForMove, 0, -2))
        if "number" == type(step) then
        else 
            step = 1
        end
        self.CursorScreenPosition.Y = self.CursorScreenPosition.Y + step
    end

    if true == isMove then
        self.KeyPressStrForMove = ""
        self:renewCursor()

    end

    self:DrawPlanets()
end

-- 计算星球所在屏幕的位置
-- rectangle 为指定宇宙位置
function _Radar.renewCursor(self)
    self.CursorScreenPosition.X, self.CursorScreenPosition.Y =
    NodeRadar:SetCursor(self.CursorScreenPosition.X, self.CursorScreenPosition.Y)
end


-- 更新 Radar 的 ScreenPlanets
-- ScreenPlanets 新的屏幕上需要显示的 planets
-- rectangle 屏幕上显示宇宙位置区域
function _Radar.RefreshScreenPlanets(self, planets, rectangle)
    local startPosition = {
        X = rectangle.Min.X,
        Y = rectangle.Min.Y
    }
    self.ScreenPlanets = {}
    for _, planet in pairs(planets) do
        -- 计算星球所在屏幕的位置
        -- rectangle 为指定宇宙位置
        planet.ScreenPosition = {
            X = planet.Info.Position.X - startPosition.X,
            Y = planet.Info.Position.Y - startPosition.Y
        }
        self.ScreenPlanets[PointToStr(planet.ScreenPosition)] = planet
    end
end

-- 画指定区域内的的星球
function _Radar.DrawPlanets(self)
    self:RefreshParInfo()
    for _, planet in pairs(self.ScreenPlanets) do
        NodeRadar:CanvasSet(
        planet.ScreenPosition.X,
        planet.ScreenPosition.Y,
        planet.Info.Character, planet.Info.ColorFg, planet.ColorBg)
    end
end

-- 画飞船
function _Radar.DrawSpaceship(self)
    self:RefreshParInfo()
    NodeRadar:CanvasSet(
    GUserSpaceship.ScreenPosition.X,
    GUserSpaceship.ScreenPosition.Y,
    GUserSpaceship.Info.Character, GUserSpaceship.Info.ColorFg, GUserSpaceship.ColorBg)
    GUserSpaceship:RefreshNodeParGUserSpaceshipStatus()
end
