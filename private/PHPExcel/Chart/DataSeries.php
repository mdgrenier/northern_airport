<?php
/**
 * PHPExcel
 *
 * Copyright (c) 2006 - 2014 PHPExcel
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * @category	PHPExcel
 * @package		PHPExcel_Chart
 * @copyright	Copyright (c) 2006 - 2014 PHPExcel (http://www.codeplex.com/PHPExcel)
 * @license		http://www.gnu.org/licenses/old-licenses/lgpl-2.1.txt	LGPL
 * @version	1.8.0, 2014-03-02
 */


/**
 * PHPExcel_Chart_DataSeries
 *
 * @category	PHPExcel
 * @package		PHPExcel_Chart
 * @copyright	Copyright (c) 2006 - 2014 PHPExcel (http://www.codeplex.com/PHPExcel)
 */
class PHPExcel_Chart_DataSeries
{

	const TYPE_BARCHART			= 'barChart';
	const TYPE_BARCHART_3D		= 'bar3DChart';
	const TYPE_LINECHART		= 'lineChart';
	const TYPE_LINECHART_3D		= 'line3DChart';
	const TYPE_AREACHART		= 'areaChart';
	const TYPE_AREACHART_3D		= 'area3DChart';
	const TYPE_PIECHART			= 'pieChart';
	const TYPE_PIECHART_3D		= 'pie3DChart';
	const TYPE_DOUGHTNUTCHART	= 'doughnutChart';
	const TYPE_DONUTCHART		= self::TYPE_DOUGHTNUTCHART;	//	Synonym
	const TYPE_SCATTERCHART		= 'scatterChart';
	const TYPE_SURFACECHART		= 'surfaceChart';
	const TYPE_SURFACECHART_3D	= 'surface3DChart';
	const TYPE_RADARCHART		= 'radarChart';
	const TYPE_BUBBLECHART		= 'bubbleChart';
	const TYPE_STOCKCHART		= 'stockChart';
	const TYPE_CANDLECHART		= self::TYPE_STOCKCHART;	   //	Synonym

	const GROUPING_CLUSTERED			= 'clustered';
	const GROUPING_STACKED				= 'stacked';
	const GROUPING_PERCENT_STACKED		= 'percentStacked';
	const GROUPING_STANDARD				= 'standard';

	const DIRECTION_BAR			= 'bar';
	const DIRECTION_HORIZONTAL	= self::DIRECTION_BAR;
	const DIRECTION_COL			= 'col';
	const DIRECTION_COLUMN		= self::DIRECTION_COL;
	const DIRECTION_VERTICAL	= self::DIRECTION_COL;

	const STYLE_LINEMARKER		= 'lineMarker';
	const STYLE_SMOOTHMARKER	= 'smoothMarker';
	const STYLE_MARKER			= 'marker';
	const STYLE_FILLED			= 'filled';


	/**
	 * Series Plot Type
	 *
	 * @var string
	 */
	private $_plotType = null;

	/**
	 * Plot Grouping Type
	 *
	 * @var boolean
	 */
	private $_plotGrouping = null;

	/**
	 * Plot Direction
	 *
	 * @var boolean
	 */
	private $_plotDirection = null;

	/**
	 * Plot Style
	 *
	 * @var string
	 */
	private $_plotStyle = null;

	/**
	 * Or