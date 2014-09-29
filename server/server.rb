require 'sinatra'
require 'json'

set :static, true
set :threaded, true
set :public_folder, 'public'

get "/" do
    redirect '/index.html'
end
get '/media/:path.js' do
    #check if file exists
    #set duration to -1 by default
    vinfo = JSON.parse(`ffprobe -loglevel error -show_format -show_streams media/#{params[:path]} -print_format json`)
    return {:duration => vinfo["format"]["duration"].to_i+1}.to_json
end

get '/media/:path.ogv' do
    #check if file exists
    $stdout.sync = true
	start= params[:start].to_f
    content_type :ogg
	status 200
    headers \
        "Access-Control-Allow-Origin" => "*",
        "Content-Disposition" => "inline", 
        "Content-Transfer-Enconding" => "binary"
    stream do |out|
        IO.popen("ffmpeg -loglevel panic -i media/#{params[:path]} -ss #{start} -f ogg -acodec libvorbis -qscale:v 10 pipe:1") do |f|
            until f.eof?
                out << f.gets(512)
            end
        end
    end
end

