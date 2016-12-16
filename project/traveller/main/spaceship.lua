local json = require("json")

local _Spaceship = {}
local _mtSpaceship = {__index = _Spaceship}

function GetSpaceshipFromDB(spaceshipId)
    sql = string.format([[
    select data from b_spaceship where spaceship_id=%d limit 1
    ]], spaceshipId)
    local rows = DB:Query(sql)
    local row = rows:FetchOne()
    rows:Close()
    if "table" ~= type(row) then
        return nil
    end

    local spaceship = NewSpaceship()
    spaceship:Format(json.decode(row.data))
    spaceship.Info.SpaceshipId = spaceshipId

    return spaceship
end

function NewSpaceshipInfo()
    return {
        SpaceshipId = nil,
        Name        = "鹦鹉螺号",
        Position    = {X = 0.0, Y = 0.0},
        Speed       = {X = 0.0, Y = 0.0},
        Character   = "x",
        ColorFg     = "blue",
        StartAt     = TimeNow(),
        Life        = 82,
        Fuel        = 100,
        Missiles    = 12,
        Jumpers     = 6,
    }
end

function NewSpaceship()
    local Spaceship           = setmetatable({}, _mtSpaceship)
    Spaceship.ScreenPosition  = {X=GetIntPart(NodeRadar:Width()/2), Y=GetIntPart(NodeRadar:Height()/2)}
    Spaceship.CenterRectangle = {}
    Spaceship:Format(NewSpaceshipInfo())
    Spaceship.ColorBg   = ""
    local Warehouse     = {}
    Spaceship.Warehouse = Warehouse

    Spaceship.lastFlushToDBForRunOneStepAt = 0

    return Spaceship
end

function _Spaceship.SetSpeedX(self, speedx)
    self:UpdateFuel(-1 * GetIntPart(math.abs(self.Info.Speed.X - speedx)))
    self.Info.Speed.X = speedx
    self:FlushToDB()
    self:RefreshNodeParGUserSpaceshipStatus()
end

function _Spaceship.SetSpeedY(self, speedy)
    self:UpdateFuel(-1 * GetIntPart(math.abs(self.Info.Speed.Y - speedy)))
    self.Info.Speed.Y = speedy
    self:FlushToDB()
    self:RefreshNodeParGUserSpaceshipStatus()
end

function _Spaceship.SetName(self, name)
    self.Info.Name = name
    self:FlushToDB()
    NodeParGUserSpaceshipStatus:SetAttribute("borderlabel", " " .. self.Info.Name .. " ")
    self:RefreshNodeParGUserSpaceshipStatus()
    NodeRadar:SetActive()
end

function _Spaceship.RefreshNodeParGUserSpaceshipStatus(self)
    NodeParGUserSpaceshipStatus:SetText(string.format([[
X: %f
Y: %f
速度X: %f/s
速度Y: %f/s
飞行历时: %d

仓库:
导弹: %d
时空跳跃者: %d
    ]], self.Info.Position.X, self.Info.Position.Y, self.Info.Speed.X, self.Info.Speed.Y, TimeNow() - self.Info.StartAt,
    self.Info.Missiles, self.Info.Jumpers))
end

-- 刷新飞船为中心的指定大小区域所在的宇宙位置
function _Spaceship.refreshCenterRectangle(self, rectangleWidth, rectangleHeight)
    local x = GetIntPart(self.Info.Position.X)
    local y = GetIntPart(self.Info.Position.Y)
    local rectangle = {}
    rectangle.Min = {
        X = GetIntPart(x - rectangleWidth / 2),
        Y = GetIntPart(y - rectangleHeight / 2)
    }
    rectangle.Max = {
        X = rectangle.Min.X + rectangleWidth,
        Y = rectangle.Min.Y + rectangleHeight,
    }

    self.CenterRectangle = rectangle

    return rectangle
end

function _Spaceship.Format(self, spaceshipInfo)
    self.Info       = {
        SpaceshipId = spaceshipInfo.SpaceshipId,
        Name        = spaceshipInfo.Name,
        Position    = spaceshipInfo.Position,
        Speed       = spaceshipInfo.Speed,
        Character   = spaceshipInfo.Character,
        ColorFg     = spaceshipInfo.ColorFg,
        StartAt     = spaceshipInfo.StartAt,
        Fuel        = spaceshipInfo.Fuel,
        Life        = spaceshipInfo.Life,
        Missiles    = spaceshipInfo.Missiles,
        Jumpers     = spaceshipInfo.Jumpers,
    }
    self:refreshCenterRectangle(NodeRadar:Width(), NodeRadar:Height())
end

function _Spaceship.FlushToDB(self)
    if "number" ~= type(self.Info.SpaceshipId) then
        return nil
    end

    sql = string.format([[
    update b_spaceship set data = '%s' where spaceship_id=%d
    ]], DB:QuoteSQL(json.encode(self.Info)), self.Info.SpaceshipId)
    local queryRet = DB:Exec(sql)
    if "string" == type(queryRet) then
        Log(queryRet)
    end
end

-- 飞船飞行，改变 position
function _Spaceship.RunOneStep(self)
    self.Info.Position.X = self.Info.Position.X + self.Info.Speed.X
    self.Info.Position.Y = self.Info.Position.Y + self.Info.Speed.Y
    self:refreshCenterRectangle(NodeRadar:Width(), NodeRadar:Height())

    if TimeNow() - self.lastFlushToDBForRunOneStepAt > 3 then
        self:FlushToDB()
    end
    self.lastFlushToDBForRunOneStepAt = TimeNow()
end

function _Spaceship.UpdateFuel(self, number)
    self.Info.Fuel = self.Info.Fuel + number
    if self.Info.Fuel < 0 then
        self.Info.Fuel = 0
    end
    NodeGaugeFuel:SetAttribute("percent", tostring(self.Info.Fuel))
    if self.Info.Fuel < 20 then
        NodeGaugeFuel:SetAttribute("percentcolor_highlighted", "black")
        NodeGaugeFuel:SetAttribute("barcolor", "red")
    elseif self.Info.Fuel < 80 then
        NodeGaugeFuel:SetAttribute("percentcolor_highlighted", "blue")
        NodeGaugeFuel:SetAttribute("barcolor", "yellow")
    else
        NodeGaugeFuel:SetAttribute("percentcolor_highlighted", "white")
        NodeGaugeFuel:SetAttribute("barcolor", "blue")
    end
    self:FlushToDB()
end

function _Spaceship.UpdateLife(self, number)
    self.Info.Life = self.Info.Life + number
    if self.Info.Life < 0 then
        self.Info.Life = 0
    end
    NodeGaugeLife:SetAttribute("percent", tostring(self.Info.Life))
    if self.Info.Life < 20 then
        NodeGaugeLife:SetAttribute("percentcolor_highlighted", "black")
        NodeGaugeLife:SetAttribute("barcolor", "red")
    elseif self.Info.Life < 80 then
        NodeGaugeLife:SetAttribute("percentcolor_highlighted", "blue")
        NodeGaugeLife:SetAttribute("barcolor", "yellow")
    else
        NodeGaugeLife:SetAttribute("percentcolor_highlighted", "black")
        NodeGaugeLife:SetAttribute("barcolor", "green")
    end
    self:FlushToDB()
end

function _Spaceship.Landing(self, number)
    GWorld:Stop(function()
        NodeModalPlanet:ModalShow()
    end)
end