package.path = package.path .. ";" .. base.ResBaseDir .. "/lua/script/?.lua;"

function Log(...)
    base.Log(unpack(arg))
end

local RandomSeed = tonumber(tostring(os.time()):reverse():sub(1, 6))

function RefreshRandomSeed()
    RandomSeed = RandomSeed + 1
    math.randomseed(RandomSeed)
end

function SetInterval(tm, callback)
    return base.SetInterval(tm, callback)
end

function SetTimeout(tm, callback)
    return base.SetTimeout(tm, callback)
end

function SendCancelSig(sig)
    base.SendCancelSig(sig)
end

function TimeNow()
    return os.time()
end

function StringSplit(str, delimiter)
    if str == nil or str == '' or delimiter == nil then
        return nil
    end

    local result = {}
    for match in (str..delimiter):gmatch("(.-)"..delimiter) do
        table.insert(result, match)
    end
    return result
end

function dumpTable(table, deepth, result)
    local i = 0

    if "table" == type(table) then
        for k, v in pairs(table) do
            i = 0; while i < deepth do result = result .. "  "; i=i+1; end

            result = result .. "[" .. tostring(k) .. "] => "
            if "table" == type(v) then
                result = result .. "{\n"
                result = dumpTable(v, deepth+1, result)
                i = 0; while i < deepth do result = result .. "  "; i=i+1; end
                result = result .. "}"
            else
                result = result .. tostring(v)
            end

            result = result .. "\n"
        end

    else
        i = 0; while i < deepth do result = result .. "  "; i=i+1; end

        result = result .. tostring(table) .. "\n"
    end

    return result
end

function DumpTable(table)
    return dumpTable(table, 0, "\n")
end

function InitRandomPoints(max, rectangle)
    local point, points = {}, {}
    local i = 0
    local isConflict = false
    while i < max do
        isConflict = false
        point = {
            X = math.random(rectangle.Min.X, rectangle.Max.X), 
            Y = math.random(rectangle.Min.Y, rectangle.Max.Y)
        }

        for _, v in pairs(points) do
            if point.X == v.X and point.Y == v.Y then
                isConflict = true
                break
            end
        end

        if false == isConflict then
            i = i + 1
            table.insert(points, point)
        end
    end

    return points
end

function CheckTableHasKey(table, key)
    for k,_ in pairs(table) do
        if k == key then
            return true
        end
    end
    return false
end

function TableLength(T)
    local count = 0
    for _ in pairs(T) do count = count + 1 end
    return count
end

function GetIntPart(x)
    if math.ceil(x) == x then
        return math.ceil(x)
    else 
        return math.ceil(x) - 1
    end
end

function PointToStr(point)
    return tostring(point.X) .. ":" .. tostring(point.Y)
end
